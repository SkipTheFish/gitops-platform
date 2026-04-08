package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
)

// DeploymentRecordStore 只负责 deployment_record 表的数据库读写。
type DeploymentRecordStore struct {
	db *gorm.DB
}

func NewDeploymentRecordStore(db *gorm.DB) *DeploymentRecordStore {
	return &DeploymentRecordStore{db: db}
}

// Create 新增一条部署记录
func (s *DeploymentRecordStore) Create(ctx context.Context, record *model.DeploymentRecord) error {
	return s.db.WithContext(ctx).Create(record).Error
}

// GetByID 按主键查询单条部署记录
func (s *DeploymentRecordStore) GetByID(ctx context.Context, id int64) (*model.DeploymentRecord, error) {
	var record model.DeploymentRecord
	err := s.db.WithContext(ctx).First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// ListByAppID 查询某个 app 的部署历史
func (s *DeploymentRecordStore) ListByAppID(ctx context.Context, appID int64) ([]model.DeploymentRecord, error) {
	var list []model.DeploymentRecord
	err := s.db.WithContext(ctx).
		Where("app_id = ?", appID).
		Order("deployed_at DESC").
		Find(&list).Error
	return list, err
}

// ListByEnvID 查询某个 environment 的部署历史
func (s *DeploymentRecordStore) ListByEnvID(ctx context.Context, envID int64) ([]model.DeploymentRecord, error) {
	var list []model.DeploymentRecord
	err := s.db.WithContext(ctx).
		Where("env_id = ?", envID).
		Order("deployed_at DESC").
		Find(&list).Error
	return list, err
}

// GetLatestByEnvID 查询某个环境最近一次部署。
// 后面做回滚时会很有用。
func (s *DeploymentRecordStore) GetLatestByEnvID(ctx context.Context, envID int64) (*model.DeploymentRecord, error) {
	var record model.DeploymentRecord
	err := s.db.WithContext(ctx).
		Where("env_id = ?", envID).
		Order("deployed_at DESC").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}
