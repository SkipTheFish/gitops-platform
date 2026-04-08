package model

import "time"

// PipelineRun 表示一次“发布流水线执行过程”。
// 这张表记录的是任务执行过程，而不是最终部署结果。
type PipelineRun struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 对应哪个应用
	AppID int64 `gorm:"column:app_id;not null;index" json:"app_id"`

	// 对应哪个环境
	EnvID int64 `gorm:"column:env_id;not null;index" json:"env_id"`

	// Git commit
	GitCommit string `gorm:"column:git_commit;type:varchar(100)" json:"git_commit"`

	// 分支名，例如 main / develop
	Branch string `gorm:"column:branch;type:varchar(100)" json:"branch"`

	// 构建出来的镜像 tag
	ImageTag string `gorm:"column:image_tag;type:varchar(255)" json:"image_tag"`

	// 当前状态：pending / running / success / failed
	Status string `gorm:"column:status;type:varchar(50);not null" json:"status"`

	// 触发方式：manual / webhook
	TriggerType string `gorm:"column:trigger_type;type:varchar(50);not null" json:"trigger_type"`

	// 开始时间
	StartedAt *time.Time `gorm:"column:started_at" json:"started_at,omitempty"`

	// 结束时间
	FinishedAt *time.Time `gorm:"column:finished_at" json:"finished_at,omitempty"`

	// 日志链接，当前先留一个字符串，后面接 CI 平台再补真实链接
	LogURL string `gorm:"column:log_url;type:text" json:"log_url"`
}

func (PipelineRun) TableName() string {
	return "pipeline_run"
}
