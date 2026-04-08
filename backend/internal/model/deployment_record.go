package model

import "time"

// DeploymentRecord 表示一次部署记录。
// 注意：这不是“当前状态”，而是“历史记录”。
// 也就是说，每部署一次，都应该新增一条记录，而不是覆盖旧记录。
type DeploymentRecord struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 这条部署记录属于哪个应用
	AppID int64 `gorm:"column:app_id;not null;index" json:"app_id"`

	// 这条部署记录属于哪个环境
	EnvID int64 `gorm:"column:env_id;not null;index" json:"env_id"`

	// 版本号，例如 v1.0.0
	Version string `gorm:"column:version;type:varchar(100);not null" json:"version"`

	// 镜像 tag，例如 order-service:v1.0.0
	ImageTag string `gorm:"column:image_tag;type:varchar(255)" json:"image_tag"`

	// 对应的 Git commit
	GitCommit string `gorm:"column:git_commit;type:varchar(100)" json:"git_commit"`

	// Argo CD 里的应用名，后面接 Argo CD API 时会用到
	ArgoCDAppName string `gorm:"column:argocd_app_name;type:varchar(255)" json:"argocd_app_name"`

	// GitOps 同步状态，例如 Synced / OutOfSync
	SyncStatus string `gorm:"column:sync_status;type:varchar(50)" json:"sync_status"`

	// 健康状态，例如 Healthy / Degraded
	HealthStatus string `gorm:"column:health_status;type:varchar(50)" json:"health_status"`

	// 操作人，后面可以接登录用户
	Operator string `gorm:"column:operator;type:varchar(100)" json:"operator"`

	// 部署时间
	DeployedAt time.Time `gorm:"column:deployed_at;autoCreateTime" json:"deployed_at"`

	// 如果这次是“回滚动作”，这里记录它是从哪个版本回滚来的。
	// 指针类型表示这个字段可以为空。
	RollbackFromVersion *string `gorm:"column:rollback_from_version;type:varchar(100)" json:"rollback_from_version,omitempty"`

	// 配置仓库本次提交的 sha
	ConfigCommitSHA string `gorm:"column:config_commit_sha;type:varchar(100)" json:"config_commit_sha"`
}

// TableName 显式指定表名，避免 GORM 猜测。
func (DeploymentRecord) TableName() string {
	return "deployment_record"
}
