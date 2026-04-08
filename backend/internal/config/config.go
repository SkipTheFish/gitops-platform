// Package config 提供应用配置的加载与管理
// 支持从.env文件和环境变量读取配置项，并提供数据库连接字符串生成方法

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv" // 用于加载 .env 文件中的环境变量
)

// 应用程序的全局配置项
// 所有字段均可通过环境变量覆盖，默认值适用于本地开发环境。
type Config struct {
	// 应用基本信息
	AppName  string
	AppEnv   string
	HTTPPort string

	// PostgreSql 数据库配置
	PGHost     string // 数据主机地址
	PGPort     string // 数据库端口
	PGUser     string // 数据库用户名
	PGPassword string // 数据库密码
	PGDB       string // 数据库名字
	PGSSLMode  string // SSL 模式

	// Redis 配置
	RedisAddr     string // Redis 地址
	RedisPassword string // Redis 密码
	RedisDB       string // 数据库编号

}

// Load 从 .env 文件和系统环境变量中加载配置，并返回 *Config 实例。
// 若 .env 文件不存在，仅使用环境变量；若环境变量未设置，则使用内置默认值。
// 注意：该函数不会返回加载错误（如 .env 不存在），仅在无法解析配置时可能出错（当前实现无错误路径）。

func Load() (*Config, error) {

	// 加载 .env 文件
	_ = godotenv.Load()

	cdg := &Config{

		// 应用配置
		AppName:  getEnv("APP_NAME", "gitops-platform"),
		AppEnv:   getEnv("APP_ENV", "dev"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),

		// PostgreSQL 配置
		PGHost:     getEnv("PG_HOST", "127.0.0.1"),
		PGPort:     getEnv("PG_PORT", "5432"),
		PGUser:     getEnv("PG_USER", "gitops"),
		PGPassword: getEnv("PG_PASSWORD", "gitops123"),
		PGDB:       getEnv("PG_DB", "gitops_platform"),
		PGSSLMode:  getEnv("PG_SSLMODE", "disable"),

		// Redis 配置
		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),
	}

	// nil 这里就是没有错误，此函数也不会返回错误
	return cdg, nil

}

// 这是一个方法
// PostgresDSN 生成 PostgreSQL 的数据源名称（DSN），
// 格式符合 pq 或 pgx 等 Go PostgreSQL 驱动的要求。
// 示例: "host=127.0.0.1 port=5432 user=gitops password=gitops123 dbname=gitops_platform sslmode=disable"
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PGHost, c.PGPort, c.PGUser, c.PGPassword, c.PGDB, c.PGSSLMode,
	)

}

// getEnv : 从环境变量中获取值
// 如果 key 对应的环境变量存在且非空，则返回其值；否则返回 fallback 默认值
func getEnv(key, fallback string) string {

	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
