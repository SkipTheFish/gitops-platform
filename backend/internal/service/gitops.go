package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type GitOpsService struct {
	repoLocalPath string
	pushBranch    string
}

// 对 push branch的构造器
func NewGitOpsService(repoLocalPath, pushBranch string) *GitOpsService {
	if strings.TrimSpace(pushBranch) == "" {
		pushBranch = "main"
	}
	return &GitOpsService{
		repoLocalPath: repoLocalPath,
		pushBranch:    pushBranch,
	}
}

type UpdateValuesInput struct {
	ValuesFilePath string
	ImageTag       string
}

type UpdateValuesResult struct {
	CommitSHA string
}

func (s *GitOpsService) UpdateImageTagAndPush(input UpdateValuesInput) (*UpdateValuesResult, error) {
	// 检查输入路径的合法性
	if strings.TrimSpace(input.ValuesFilePath) == "" {
		return nil, fmt.Errorf("values_file_path is required")
	}
	// 检查镜像标签的合法性
	if strings.TrimSpace(input.ImageTag) == "" {
		return nil, fmt.Errorf("image_tag is required")
	}

	absPath := filepath.Join(s.repoLocalPath, input.ValuesFilePath)

	// 修改文件内容
	if err := s.updateValuesFile(absPath, input.ImageTag); err != nil {
		return nil, err
	}

	// git add
	if err := runCmd(s.repoLocalPath, "git", "add", input.ValuesFilePath); err != nil {
		return nil, err
	}

	// git commit
	commitMsg := fmt.Sprintf("chore(release): update image tag to %s", input.ImageTag)
	commitErr := runCmd(s.repoLocalPath, "git", "commit", "-m", commitMsg)
	if commitErr != nil {
		// 没变化时可视为幂等成功
		if strings.Contains(commitErr.Error(), "nothing to commit") ||
			strings.Contains(commitErr.Error(), "nothing added to commit") {
			sha, err := outputCmd(s.repoLocalPath, "git", "rev-parse", "HEAD")
			if err != nil {
				return nil, err
			}
			return &UpdateValuesResult{CommitSHA: sha}, nil
		}
		return nil, commitErr
	}

	// git push
	if err := runCmd(s.repoLocalPath, "git", "push", "origin", s.pushBranch); err != nil {
		return nil, err
	}

	// 获取最终的 Commit SHA 供审计
	sha, err := outputCmd(s.repoLocalPath, "git", "rev-parse", "HEAD")
	if err != nil {
		return nil, err
	}

	return &UpdateValuesResult{
		CommitSHA: sha,
	}, nil
}

// 文件IO与内容解析
func (s *GitOpsService) updateValuesFile(path string, imageTag string) error {
	data, err := os.ReadFile(path)

	// 检验读取文件
	if err != nil {
		return fmt.Errorf("read values file failed: %w", err)
	}

	// 反序列化，[]byte 类型的 YAML 变量转化到values当中
	// map[string]any —— 它可以表示任意层级的键值对结构：
	// {"name": value}这样，就可以根据字段名查找对应的值
	// 索引是string，值可以是any，甚至可以也是键值对
	// 这里的values["images"]就是一个值是["tag"]的键值对
	// eg: values"images" = "tag": v1这样
	var values map[string]any

	if err := yaml.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("unmarshal yaml failed: %w", err)
	}

	// 将 YAML 中"image"的值读取到imageObj当中
	imageObj, ok := values["image"].(map[string]any)
	if !ok {
		// 处理 YAML 中image格式错误的情况
		imageObj = map[string]any{} // 分配一块新的内存空间，变量类型是map格式
		values["image"] = imageObj  // 将这个新的map节点与values["image"]绑定
	}
	// 无论images是成功从values里面取到的旧map
	// 还是格式错误产生的新 map节点
	// 都指向同一个values["image"]
	imageObj["tag"] = imageTag

	// 序列化返回YAML文件
	newData, err := yaml.Marshal(values)
	if err != nil {
		return fmt.Errorf("marshal yaml failed: %w", err)
	}

	// 覆盖写入
	if err := os.WriteFile(path, newData, 0644); err != nil {
		return fmt.Errorf("write values file failed: %w", err)
	}
	return nil
}

func runCmd(dir string, name string, args ...string) error {
	// 定义日志器
	cmd := exec.Command(name, args...)
	// cd
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %v failed: %v, output=%s", name, args, err, string(output))
	}
	return nil
}

func outputCmd(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %v failed: %v, output=%s", name, args, err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}
