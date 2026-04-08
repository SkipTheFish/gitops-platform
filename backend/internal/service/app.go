package service

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
	"github.com/SkipTheFish/gitops-platform/backend/internal/store"
)

type AppService struct {
	appStore *store.AppStore
}

func NewAppService(appStore *store.AppStore) *AppService {
	return &AppService{
		appStore: appStore,
	}
}

// 读取App_handler请求体的数据读取到input
type CreateAppInput struct {
	Name           string
	RepoURL        string
	ConfigRepoURL  string
	ClusterName    string
	Namespace      string
	HelmChartPath  string
	ValuesFilePath string
}

type UpdateAppInput struct {
	ID             int64
	Name           string
	RepoURL        string
	ConfigRepoURL  string
	ClusterName    string
	Namespace      string
	HelmChartPath  string
	ValuesFilePath string
}

// Create 实现“新增应用”的核心逻辑。
//
// 思路：
// 1. 先做基础清洗（trim space）
// 2. 做业务校验（必填、名称唯一）
// 3. 组装 model
// 4. 调 store 落库

// context.Context 是什么？
// 在 Go 语言中，Context（上下文）被形象地称为**“请求的身份证”或“遥控器”**。
// 它不是前端传来的文本，也不负责把对象转成 JSON（那是序列化）。
// 它的核心作用是：告知程序什么时候该“停下来”。
// 想象一个实际场景：
// 用户在网页上点击“查询应用详情”。
// 你的后端接收到请求，开始执行 GetByID。
// 突然，用户觉得等得太久，直接把浏览器标签页关掉了。

// 如果没有 ctx：
// 你的 AppService 会继续傻傻地调用 AppStore，AppStore 会继续去数据库执行 SQL。数据库查完后返回给 Service，Service 再返回给 Handler。最后 Handler 准备发给前端时，发现：咦？人呢？连接断了！ > 结果：白忙活一场，浪费了 CPU 和数据库连接资源。

// 如果有 ctx：
// 用户关掉网页的瞬间，Gin 框架会立刻往 ctx 里发送一个**“取消信号”**。

// AppStore 里的 s.db.WithContext(ctx) 接收到了这个信号。
// GORM 发现信号后，会立即中断正在进行的 SQL 查询并报错退出。
// 整个调用链路瞬间停止，服务器资源立即释放。
func (s *AppService) Create(ctx context.Context, input CreateAppInput) (*model.App, error) {

	// 清晰input里面的数据
	input.Name = strings.TrimSpace(input.Name)
	input.RepoURL = strings.TrimSpace(input.RepoURL)
	input.ConfigRepoURL = strings.TrimSpace(input.ConfigRepoURL)
	input.ClusterName = strings.TrimSpace(input.ClusterName)
	input.Namespace = strings.TrimSpace(input.Namespace)
	input.HelmChartPath = strings.TrimSpace(input.HelmChartPath)
	input.ValuesFilePath = strings.TrimSpace(input.ValuesFilePath)

	// 检查有无空置
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	if input.RepoURL == "" {
		return nil, errors.New("repo_url is required")
	}
	if input.ConfigRepoURL == "" {
		return nil, errors.New("config_repo_url is required")
	}
	if input.ClusterName == "" {
		return nil, errors.New("cluster_name is required")
	}
	if input.Namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if input.HelmChartPath == "" {
		return nil, errors.New("helm_chart_path is required")
	}
	if input.ValuesFilePath == "" {
		return nil, errors.New("values_file_path is required")
	}

	// 名称唯一性校验，调用appStore里的方法
	_, err := s.appStore.GetByName(ctx, input.Name)
	if err == nil {
		// 如果没报错，说明查到了同名的，这在业务上是不允许的
		return nil, errors.New("app name already exists")
	}
	// 查找不到
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 一切ok，创建一个appModel 实体对象，接下来store会保存这个model
	app := &model.App{
		Name:           input.Name,
		RepoURL:        input.RepoURL,
		ConfigRepoURL:  input.ConfigRepoURL,
		ClusterName:    input.ClusterName,
		Namespace:      input.Namespace,
		HelmChartPath:  input.HelmChartPath,
		ValuesFilePath: input.ValuesFilePath,
	}

	// 	这时候，app 这个对象已经是一个完整的、合法的业务实体了。
	// Service 告诉 Store：“我检查过了，没问题，你把它存进数据库吧。”
	if err := s.appStore.Create(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

// 查询应用列表
func (s *AppService) List(ctx context.Context) ([]model.App, error) {
	return s.appStore.List(ctx)
}

// 查询单个应用详情（获取主键id）
func (s *AppService) GetByID(ctx context.Context, id int64) (*model.App, error) {
	if id <= 0 {
		return nil, errors.New("invalid app id")
	}
	return s.appStore.GetByID(ctx, id)
}

// Update 的思想：
// 1. 先查旧数据
// 2. 只更新用户真正传了的新值
// 3. 如果要改 name，还要做唯一性校验
func (s *AppService) Update(ctx context.Context, input UpdateAppInput) (*model.App, error) {
	// 检查输入的app id
	if input.ID <= 0 {
		return nil, errors.New("invalid app id")
	}

	//	查找app id，把旧的数据存储到app当中
	app, err := s.appStore.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// 输入的name不为空并且新的输入和旧的名字不一样
	if strings.TrimSpace(input.Name) != "" && strings.TrimSpace(input.Name) != app.Name {
		// 赋值
		targetName := strings.TrimSpace(input.Name)

		// 用输入的名字来查询id，得到的数据存储到existed中
		existed, err := s.appStore.GetByName(ctx, targetName)

		// 如果existed.id 和 旧数据id相等，说明名字已经存在了
		if err == nil && existed.ID != app.ID {
			return nil, errors.New("app name already exists")
		}

		// 除了“没找到记录”之外的任何数据库错误，都要拦截并报错。（系统级异常）
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// 一切正常，替换名字
		app.Name = targetName
	}
	// 替换其他数据
	if strings.TrimSpace(input.RepoURL) != "" {
		app.RepoURL = strings.TrimSpace(input.RepoURL)
	}
	if strings.TrimSpace(input.ConfigRepoURL) != "" {
		app.ConfigRepoURL = strings.TrimSpace(input.ConfigRepoURL)
	}
	if strings.TrimSpace(input.ClusterName) != "" {
		app.ClusterName = strings.TrimSpace(input.ClusterName)
	}
	if strings.TrimSpace(input.Namespace) != "" {
		app.Namespace = strings.TrimSpace(input.Namespace)
	}
	if strings.TrimSpace(input.HelmChartPath) != "" {
		app.HelmChartPath = strings.TrimSpace(input.HelmChartPath)
	}
	if strings.TrimSpace(input.ValuesFilePath) != "" {
		app.ValuesFilePath = strings.TrimSpace(input.ValuesFilePath)
	}

	// 调用store层，把数据更新到数据库
	if err := s.appStore.Update(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}
