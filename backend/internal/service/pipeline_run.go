package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
	"github.com/SkipTheFish/gitops-platform/backend/internal/store"
)

// PipelineRunService 负责“手动触发发布”和“模拟流水线执行”。
type PipelineRunService struct {
	pipelineStore   *store.PipelineRunStore
	appStore        *store.AppStore
	envStore        *store.EnvironmentStore
	deploymentStore *store.DeploymentRecordStore
	auditStore      *store.OperationAuditStore
	executor        *PipelineExecutor
}

func NewPipelineRunService(
	pipelineStore *store.PipelineRunStore,
	appStore *store.AppStore,
	envStore *store.EnvironmentStore,
	deploymentStore *store.DeploymentRecordStore,
	auditStore *store.OperationAuditStore,
	executor *PipelineExecutor,
) *PipelineRunService {
	return &PipelineRunService{
		pipelineStore:   pipelineStore,
		appStore:        appStore,
		envStore:        envStore,
		deploymentStore: deploymentStore,
		auditStore:      auditStore,
		executor:        executor,
	}
}

// CreatePipelineRunInput 表示“手动触发发布”时需要的输入。
type CreatePipelineRunInput struct {
	AppID       int64
	EnvID       int64
	GitCommit   string
	Branch      string
	ImageTag    string
	TriggerType string
	Operator    string
	Version     string
}

// CreateManualRun 的思路：
// 1. 校验 app / env 是否存在
// 2. 校验 env 属于 app
// 3. 创建一条 pending 状态的 pipeline_run
// 4. 写一条 audit
// 5. 启动后台 goroutine 模拟执行
func (s *PipelineRunService) CreateManualRun(ctx context.Context, input CreatePipelineRunInput) (*model.PipelineRun, error) {
	if input.AppID <= 0 {
		return nil, errors.New("invalid app_id")
	}
	if input.EnvID <= 0 {
		return nil, errors.New("invalid env_id")
	}

	input.GitCommit = strings.TrimSpace(input.GitCommit)
	input.Branch = strings.TrimSpace(input.Branch)
	input.ImageTag = strings.TrimSpace(input.ImageTag)
	input.TriggerType = strings.TrimSpace(input.TriggerType)
	input.Operator = strings.TrimSpace(input.Operator)
	input.Version = strings.TrimSpace(input.Version)

	if input.Branch == "" {
		input.Branch = "main"
	}
	if input.TriggerType == "" {
		input.TriggerType = "manual"
	}
	if input.Operator == "" {
		return nil, errors.New("operator is required")
	}
	if input.Version == "" {
		return nil, errors.New("version is required")
	}

	app, err := s.appStore.GetByID(ctx, input.AppID)
	if err != nil {
		return nil, errors.New("app not found")
	}

	env, err := s.envStore.GetByID(ctx, input.EnvID)
	if err != nil {
		return nil, errors.New("environment not found")
	}

	if env.AppID != app.ID {
		return nil, errors.New("environment does not belong to this app")
	}

	run := &model.PipelineRun{
		AppID:       input.AppID,
		EnvID:       input.EnvID,
		GitCommit:   input.GitCommit,
		Branch:      input.Branch,
		ImageTag:    input.ImageTag,
		Status:      "pending",
		TriggerType: input.TriggerType,
		LogURL:      "",
	}

	if err := s.pipelineStore.Create(ctx, run); err != nil {
		return nil, err
	}

	// 记录一条“触发发布”的审计
	_ = s.auditStore.Create(ctx, &model.OperationAudit{
		Operator:   input.Operator,
		ActionType: "trigger_pipeline",
		TargetID:   run.ID,
		Detail:     fmt.Sprintf("trigger pipeline for app=%s env=%s version=%s", app.Name, env.EnvName, input.Version),
	})

	// 异步模拟流水线执行
	go s.executor.Execute(run.ID, input)

	return run, nil
}

// simulatePipeline 模拟一个本地流水线。
// 这是第 4 步的关键：先不接真实 Jenkins，而是把“状态流”跑通。
// func (s *PipelineRunService) simulatePipeline(runID int64, input CreatePipelineRunInput) {
// 	ctx := context.Background()

// 	// 1. 先把状态改成 running
// 	run, err := s.pipelineStore.GetByID(ctx, runID)
// 	if err != nil {
// 		return
// 	}

// 	now := time.Now()
// 	run.Status = "running"
// 	run.StartedAt = &now
// 	run.LogURL = fmt.Sprintf("/fake-logs/pipeline/%d", run.ID)
// 	if err := s.pipelineStore.Update(ctx, run); err != nil {
// 		return
// 	}

// 	// 2. 模拟执行过程
// 	time.Sleep(2 * time.Second) // 模拟构建镜像
// 	time.Sleep(2 * time.Second) // 模拟推送镜像
// 	time.Sleep(2 * time.Second) // 模拟更新配置 / 触发部署

// 	// 3. 默认模拟成功
// 	finishedAt := time.Now()
// 	run.Status = "success"
// 	run.FinishedAt = &finishedAt
// 	if err := s.pipelineStore.Update(ctx, run); err != nil {
// 		return
// 	}

// 	// 4. 流水线成功后，自动写 deployment_record
// 	// 这样前端在“部署记录”页面就能立刻看到结果。
// 	deployRecord := &model.DeploymentRecord{
// 		AppID:         input.AppID,
// 		EnvID:         input.EnvID,
// 		Version:       input.Version,
// 		ImageTag:      input.ImageTag,
// 		GitCommit:     input.GitCommit,
// 		ArgoCDAppName: fmt.Sprintf("app-%d-env-%d", input.AppID, input.EnvID),
// 		SyncStatus:    "Synced",
// 		HealthStatus:  "Healthy",
// 		Operator:      input.Operator,
// 	}

// 	if err := s.deploymentStore.Create(ctx, deployRecord); err != nil {
// 		return
// 	}

// 	// 5. 记录部署成功的审计日志
// 	_ = s.auditStore.Create(ctx, &model.OperationAudit{
// 		Operator:   input.Operator,
// 		ActionType: "deploy",
// 		TargetID:   deployRecord.ID,
// 		Detail:     fmt.Sprintf("deployment success, version=%s image=%s", input.Version, input.ImageTag),
// 	})
// }

func (s *PipelineRunService) GetByID(ctx context.Context, id int64) (*model.PipelineRun, error) {
	if id <= 0 {
		return nil, errors.New("invalid pipeline id")
	}
	return s.pipelineStore.GetByID(ctx, id)
}

func (s *PipelineRunService) ListByAppID(ctx context.Context, appID int64) ([]model.PipelineRun, error) {
	if appID <= 0 {
		return nil, errors.New("invalid app id")
	}
	if _, err := s.appStore.GetByID(ctx, appID); err != nil {
		return nil, errors.New("app not found")
	}
	return s.pipelineStore.ListByAppID(ctx, appID)
}

func (s *PipelineRunService) ListByEnvID(ctx context.Context, envID int64) ([]model.PipelineRun, error) {
	if envID <= 0 {
		return nil, errors.New("invalid environment id")
	}
	if _, err := s.envStore.GetByID(ctx, envID); err != nil {
		return nil, errors.New("environment not found")
	}
	return s.pipelineStore.ListByEnvID(ctx, envID)
}
