package service

import (
	"context"
	"errors"
	"strings"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
	"github.com/SkipTheFish/gitops-platform/backend/internal/store"
)

// DeploymentRecordService 负责“部署记录”的业务逻辑。
type DeploymentRecordService struct {
	deploymentStore *store.DeploymentRecordStore
	appStore        *store.AppStore
	envStore        *store.EnvironmentStore
	auditStore      *store.OperationAuditStore
}

// NewDeploymentRecordService 构造函数。
// 这里注入多个 store，是因为创建部署记录时需要联动校验 app / env，并写 audit。
func NewDeploymentRecordService(
	deploymentStore *store.DeploymentRecordStore,
	appStore *store.AppStore,
	envStore *store.EnvironmentStore,
	auditStore *store.OperationAuditStore,
) *DeploymentRecordService {
	return &DeploymentRecordService{
		deploymentStore: deploymentStore,
		appStore:        appStore,
		envStore:        envStore,
		auditStore:      auditStore,
	}
}

// CreateDeploymentRecordInput 表示 service 层创建部署记录所需的参数。
// 注意它不是 HTTP 请求体，而是“业务输入对象”。
type CreateDeploymentRecordInput struct {
	AppID               int64
	EnvID               int64
	Version             string
	ImageTag            string
	GitCommit           string
	ArgoCDAppName       string
	SyncStatus          string
	HealthStatus        string
	Operator            string
	RollbackFromVersion *string
}

// Create 的思路：
// 1. 做参数清洗和校验
// 2. 确认 app 存在
// 3. 确认 env 存在
// 4. 确认这个 env 真的属于这个 app
// 5. 写 deployment_record
// 6. 顺手写 operation_audit
func (s *DeploymentRecordService) Create(ctx context.Context, input CreateDeploymentRecordInput) (*model.DeploymentRecord, error) {
	if input.AppID <= 0 {
		return nil, errors.New("invalid app_id")
	}
	if input.EnvID <= 0 {
		return nil, errors.New("invalid env_id")
	}

	input.Version = strings.TrimSpace(input.Version)
	input.ImageTag = strings.TrimSpace(input.ImageTag)
	input.GitCommit = strings.TrimSpace(input.GitCommit)
	input.ArgoCDAppName = strings.TrimSpace(input.ArgoCDAppName)
	input.SyncStatus = strings.TrimSpace(input.SyncStatus)
	input.HealthStatus = strings.TrimSpace(input.HealthStatus)
	input.Operator = strings.TrimSpace(input.Operator)

	if input.Version == "" {
		return nil, errors.New("version is required")
	}
	if input.Operator == "" {
		return nil, errors.New("operator is required")
	}

	// 校验 app 是否存在
	app, err := s.appStore.GetByID(ctx, input.AppID)
	if err != nil {
		return nil, errors.New("app not found")
	}

	// 校验 env 是否存在
	env, err := s.envStore.GetByID(ctx, input.EnvID)
	if err != nil {
		return nil, errors.New("environment not found")
	}

	// 非常关键：校验 environment 是否属于这个 app
	if env.AppID != app.ID {
		return nil, errors.New("environment does not belong to this app")
	}

	record := &model.DeploymentRecord{
		AppID:               input.AppID,
		EnvID:               input.EnvID,
		Version:             input.Version,
		ImageTag:            input.ImageTag,
		GitCommit:           input.GitCommit,
		ArgoCDAppName:       input.ArgoCDAppName,
		SyncStatus:          input.SyncStatus,
		HealthStatus:        input.HealthStatus,
		Operator:            input.Operator,
		RollbackFromVersion: input.RollbackFromVersion,
	}

	if err := s.deploymentStore.Create(ctx, record); err != nil {
		return nil, err
	}

	// 创建一条审计日志
	// 为什么这里顺手写？
	// 因为“创建部署记录”本身就是一个重要操作，后面排查发布问题时会用到。
	audit := &model.OperationAudit{
		Operator:   input.Operator,
		ActionType: "deploy",
		TargetID:   record.ID,
		Detail:     "deploy version=" + input.Version + ", app=" + app.Name + ", env=" + env.EnvName,
	}
	_ = s.auditStore.Create(ctx, audit)

	return record, nil
}

// GetByID 查询单条部署记录
func (s *DeploymentRecordService) GetByID(ctx context.Context, id int64) (*model.DeploymentRecord, error) {
	if id <= 0 {
		return nil, errors.New("invalid deployment id")
	}
	return s.deploymentStore.GetByID(ctx, id)
}

// ListByAppID 查询某个应用的部署历史
func (s *DeploymentRecordService) ListByAppID(ctx context.Context, appID int64) ([]model.DeploymentRecord, error) {
	if appID <= 0 {
		return nil, errors.New("invalid app id")
	}

	_, err := s.appStore.GetByID(ctx, appID)
	if err != nil {
		return nil, errors.New("app not found")
	}

	return s.deploymentStore.ListByAppID(ctx, appID)
}

// ListByEnvID 查询某个环境的部署历史
func (s *DeploymentRecordService) ListByEnvID(ctx context.Context, envID int64) ([]model.DeploymentRecord, error) {
	if envID <= 0 {
		return nil, errors.New("invalid environment id")
	}

	_, err := s.envStore.GetByID(ctx, envID)
	if err != nil {
		return nil, errors.New("environment not found")
	}

	return s.deploymentStore.ListByEnvID(ctx, envID)
}
