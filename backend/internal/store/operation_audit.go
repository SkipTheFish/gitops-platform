package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
)

// OperationAuditStore 负责 operation_audit 表的读写。
type OperationAuditStore struct {
	db *gorm.DB
}

func NewOperationAuditStore(db *gorm.DB) *OperationAuditStore {
	return &OperationAuditStore{db: db}
}

// Create 新增审计记录
func (s *OperationAuditStore) Create(ctx context.Context, audit *model.OperationAudit) error {
	return s.db.WithContext(ctx).Create(audit).Error
}

// ListByTargetID 查询某个目标对象的审计日志
func (s *OperationAuditStore) ListByTargetID(ctx context.Context, targetID int64) ([]model.OperationAudit, error) {
	var list []model.OperationAudit
	err := s.db.WithContext(ctx).
		Where("target_id = ?", targetID).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}
