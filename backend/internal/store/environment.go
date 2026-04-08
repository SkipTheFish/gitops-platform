package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
)

type EnvironmentStore struct {
	db *gorm.DB
}

func NewEnvironmentStore(db *gorm.DB) *EnvironmentStore {
	return &EnvironmentStore{db: db}
}

func (s *EnvironmentStore) Create(ctx context.Context, env *model.Environment) error {
	return s.db.WithContext(ctx).Create(env).Error
}

func (s *EnvironmentStore) ListByAppID(ctx context.Context, appID int64) ([]model.Environment, error) {
	var list []model.Environment
	err := s.db.WithContext(ctx).
		Where("app_id = ?", appID).
		Order("id asc").
		Find(&list).Error
	return list, err
}

func (s *EnvironmentStore) GetByID(ctx context.Context, id int64) (*model.Environment, error) {
	var env model.Environment
	err := s.db.WithContext(ctx).First(&env, id).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}

// 用于校验同一个 app 下 env_name 是否重复
func (s *EnvironmentStore) GetByAppIDAndEnvName(ctx context.Context, appID int64, envName string) (*model.Environment, error) {
	var env model.Environment
	err := s.db.WithContext(ctx).
		Where("app_id = ? AND env_name = ?", appID, envName).
		First(&env).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (s *EnvironmentStore) Update(ctx context.Context, env *model.Environment) error {
	return s.db.WithContext(ctx).Save(env).Error
}
