// 让 Go 结构体和数据库表结构一一对应。

package model

import "time"

// App 表示“应用元信息”。
// 这张表不关心某一次发布，只关心“这个应用是谁、代码和配置仓库在哪、默认部署信息是什么”。
type App struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"column:name;type:varchar(100);not null;uniqueIndex" json:"name"`
	RepoURL        string    `gorm:"column:repo_url;type:text;not null" json:"repo_url"`
	ConfigRepoURL  string    `gorm:"column:config_repo_url;type:text;not null" json:"config_repo_url"`
	ClusterName    string    `gorm:"column:cluster_name;type:varchar(100);not null" json:"cluster_name"`
	Namespace      string    `gorm:"column:namespace;type:varchar(100);not null" json:"namespace"`
	HelmChartPath  string    `gorm:"column:helm_chart_path;type:text;not null" json:"helm_chart_path"`
	ValuesFilePath string    `gorm:"column:values_file_path;type:text;not null" json:"values_file_path"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 告诉 GORM 这个结构体对应哪张表。
// 如果不写，GORM 会自动猜表名，但工程里显式写出来更稳。

func (App) TableName() string {
	return "app"
}
