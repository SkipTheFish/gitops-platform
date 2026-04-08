// 创建一个 GORM 数据库连接对象 (*gorm.DB) 供业务逻辑使用
// 注意：只是初始化配置和驱动，没有立刻连接

package pkg

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 输入：PostgreSQL 的 DSN
// 输出：GORM 的数据库句柄（*gorm.DB）或错误
func NewPostgres(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
