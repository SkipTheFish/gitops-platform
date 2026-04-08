package service

import (
	"context"
	"errors"
	"strings"

	"github.com/SkipTheFish/gitops-platform/backend/internal/model"
	"github.com/SkipTheFish/gitops-platform/backend/internal/store"
)

// OperationAuditService 负责操作审计的业务逻辑。
type OperationAuditService struct {
	auditStore *store.OperationAuditStore
}

func NewOperationAuditService(auditStore *store.OperationAuditStore) *OperationAuditService {
	return &OperationAuditService{
		auditStore: auditStore,
	}
}

// CreateOperationAuditInput 是创建审计记录的输入参数。
type CreateOperationAuditInput struct {
	Operator   string
	ActionType string
	TargetID   int64
	Detail     string
}

// Create 创建审计记录
func (s *OperationAuditService) Create(ctx context.Context, input CreateOperationAuditInput) (*model.OperationAudit, error) {
	input.Operator = strings.TrimSpace(input.Operator)
	input.ActionType = strings.TrimSpace(input.ActionType)
	input.Detail = strings.TrimSpace(input.Detail)

	if input.Operator == "" {
		return nil, errors.New("operator is required")
	}
	if input.ActionType == "" {
		return nil, errors.New("action_type is required")
	}
	if input.TargetID <= 0 {
		return nil, errors.New("target_id is required")
	}

	audit := &model.OperationAudit{
		Operator:   input.Operator,
		ActionType: input.ActionType,
		TargetID:   input.TargetID,
		Detail:     input.Detail,
	}

	if err := s.auditStore.Create(ctx, audit); err != nil {
		return nil, err
	}
	return audit, nil
}

// ListByTargetID 查询某个目标对象的审计历史
func (s *OperationAuditService) ListByTargetID(ctx context.Context, targetID int64) ([]model.OperationAudit, error) {
	if targetID <= 0 {
		return nil, errors.New("invalid target id")
	}
	return s.auditStore.ListByTargetID(ctx, targetID)
}
