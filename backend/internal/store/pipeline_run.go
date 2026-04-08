package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
)

// PipelineRunStore 只负责 pipeline_run 表的数据库读写。
type PipelineRunStore struct {
	db *gorm.DB
}

func NewPipelineRunStore(db *gorm.DB) *PipelineRunStore {
	return &PipelineRunStore{db: db}
}

// Create 新增一条流水线记录
func (s *PipelineRunStore) Create(ctx context.Context, run *model.PipelineRun) error {
	return s.db.WithContext(ctx).Create(run).Error
}

// GetByID 查询单条流水线记录
func (s *PipelineRunStore) GetByID(ctx context.Context, id int64) (*model.PipelineRun, error) {
	var run model.PipelineRun
	err := s.db.WithContext(ctx).First(&run, id).Error
	if err != nil {
		return nil, err
	}
	return &run, nil
}

// ListByAppID 查询某个应用的流水线记录
func (s *PipelineRunStore) ListByAppID(ctx context.Context, appID int64) ([]model.PipelineRun, error) {
	var list []model.PipelineRun
	err := s.db.WithContext(ctx).
		Where("app_id = ?", appID).
		Order("id DESC").
		Find(&list).Error
	return list, err
}

// ListByEnvID 查询某个环境的流水线记录
func (s *PipelineRunStore) ListByEnvID(ctx context.Context, envID int64) ([]model.PipelineRun, error) {
	var list []model.PipelineRun
	err := s.db.WithContext(ctx).
		Where("env_id = ?", envID).
		Order("id DESC").
		Find(&list).Error
	return list, err
}

// Update 更新流水线记录
func (s *PipelineRunStore) Update(ctx context.Context, run *model.PipelineRun) error {
	return s.db.WithContext(ctx).Save(run).Error
}
