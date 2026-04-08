package service

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
	"github.com/SkipTheFish/gitops-platform/backend/internal/store"
)

type EnvironmentService struct {
	envStore *store.EnvironmentStore
	appStore *store.AppStore
}

func NewEnvironmentService(envStore *store.EnvironmentStore, appStore *store.AppStore) *EnvironmentService {
	return &EnvironmentService{
		envStore: envStore,
		appStore: appStore,
	}
}

type CreateEnvironmentInput struct {
	AppID           int64
	EnvName         string
	ClusterName     string
	Namespace       string
	AutoSyncEnabled bool
	ValuesFilePath  string
	ArgoCDAppName   string
}

type UpdateEnvironmentInput struct {
	ID              int64
	EnvName         string
	ClusterName     string
	Namespace       string
	AutoSyncEnabled *bool
	// 这个环境对应改那个 values 文件
	ValuesFilePath string

	// 这个环境在 Argo CD 中对应的 Application 名称
	ArgoCDAppName string
}

// Create 创建环境的思路：
// 1. 先确认 app 存在
// 2. 校验 env_name 在该 app 下不能重复
// 3. 保存到 environment 表
func (s *EnvironmentService) Create(ctx context.Context, input CreateEnvironmentInput) (*model.Environment, error) {
	if input.AppID <= 0 {
		return nil, errors.New("invalid app id")
	}

	input.EnvName = strings.TrimSpace(input.EnvName)
	input.ClusterName = strings.TrimSpace(input.ClusterName)
	input.Namespace = strings.TrimSpace(input.Namespace)

	if input.EnvName == "" {
		return nil, errors.New("env_name is required")
	}
	if input.ClusterName == "" {
		return nil, errors.New("cluster_name is required")
	}
	if input.Namespace == "" {
		return nil, errors.New("namespace is required")
	}

	// 先确保 app 存在
	_, err := s.appStore.GetByID(ctx, input.AppID)
	if err != nil {
		return nil, errors.New("app not found")
	}

	// 同一个 app 下 env_name 不能重复
	_, err = s.envStore.GetByAppIDAndEnvName(ctx, input.AppID, input.EnvName)
	if err == nil {
		return nil, errors.New("environment already exists for this app")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	env := &model.Environment{
		AppID:           input.AppID,
		EnvName:         input.EnvName,
		ClusterName:     input.ClusterName,
		Namespace:       input.Namespace,
		AutoSyncEnabled: input.AutoSyncEnabled,
		ValuesFilePath:  strings.TrimSpace(input.ValuesFilePath),
		ArgoCDAppName:   strings.TrimSpace(input.ArgoCDAppName),
	}

	if err := s.envStore.Create(ctx, env); err != nil {
		return nil, err
	}
	return env, nil
}

func (s *EnvironmentService) ListByAppID(ctx context.Context, appID int64) ([]model.Environment, error) {
	if appID <= 0 {
		return nil, errors.New("invalid app id")
	}

	_, err := s.appStore.GetByID(ctx, appID)
	if err != nil {
		return nil, errors.New("app not found")
	}

	return s.envStore.ListByAppID(ctx, appID)
}

func (s *EnvironmentService) GetByID(ctx context.Context, id int64) (*model.Environment, error) {
	if id <= 0 {
		return nil, errors.New("invalid environment id")
	}
	return s.envStore.GetByID(ctx, id)
}

func (s *EnvironmentService) Update(ctx context.Context, input UpdateEnvironmentInput) (*model.Environment, error) {
	if input.ID <= 0 {
		return nil, errors.New("invalid environment id")
	}

	env, err := s.envStore.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.EnvName) != "" && strings.TrimSpace(input.EnvName) != env.EnvName {
		targetEnvName := strings.TrimSpace(input.EnvName)

		existed, err := s.envStore.GetByAppIDAndEnvName(ctx, env.AppID, targetEnvName)
		if err == nil && existed.ID != env.ID {
			return nil, errors.New("environment name already exists for this app")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		env.EnvName = targetEnvName
	}

	if strings.TrimSpace(input.ClusterName) != "" {
		env.ClusterName = strings.TrimSpace(input.ClusterName)
	}
	if strings.TrimSpace(input.Namespace) != "" {
		env.Namespace = strings.TrimSpace(input.Namespace)
	}
	if input.AutoSyncEnabled != nil {
		env.AutoSyncEnabled = *input.AutoSyncEnabled
	}

	if input.ValuesFilePath != "" {
		env.ValuesFilePath = input.ValuesFilePath
	}

	if input.ArgoCDAppName != "" {
		env.ArgoCDAppName = input.ArgoCDAppName
	}

	if err := s.envStore.Update(ctx, env); err != nil {
		return nil, err
	}
	return env, nil
}
