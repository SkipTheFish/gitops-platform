package model

import "time"

// OperationAudit 表示一条操作审计日志。
// 它关注的是“谁做了什么”，不是“部署结果是什么”。
type OperationAudit struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 谁做的
	Operator string `gorm:"column:operator;type:varchar(100);not null" json:"operator"`

	// 做了什么，例如 deploy / rollback / create_app / update_env
	ActionType string `gorm:"column:action_type;type:varchar(50);not null" json:"action_type"`

	// 作用到哪个目标对象上
	TargetID int64 `gorm:"column:target_id;not null" json:"target_id"`

	// 详细信息
	Detail string `gorm:"column:detail;type:text" json:"detail"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (OperationAudit) TableName() string {
	return "operation_audit"
}
