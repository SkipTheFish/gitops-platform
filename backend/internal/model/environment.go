package model

import "time"

// Environment 表示应用在某个环境下的部署信息。
// 这里强调“同一个 app 可以有多个 environment”。
type Environment struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AppID           int64     `gorm:"column:app_id;not null;index" json:"app_id"`
	EnvName         string    `gorm:"column:env_name;type:varchar(50);not null" json:"env_name"`
	ClusterName     string    `gorm:"column:cluster_name;type:varchar(100);not null" json:"cluster_name"`
	Namespace       string    `gorm:"column:namespace;type:varchar(100);not null" json:"namespace"`
	AutoSyncEnabled bool      `gorm:"column:auto_sync_enabled;not null;default:false" json:"auto_sync_enabled"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	// 第 5 步新增：这个环境对应要改哪个 values 文件
	ValuesFilePath string `gorm:"column:values_file_path;type:varchar(255)" json:"values_file_path"`

	// 第 5 步新增：这个环境在 Argo CD 中对应的 Application 名称
	ArgoCDAppName string `gorm:"column:argocd_app_name;type:varchar(255)" json:"argocd_app_name"`
}

func (Environment) TableName() string {
	return "environment"
}
