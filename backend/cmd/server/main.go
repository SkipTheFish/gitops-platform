package main

import (
	"log"

	"github.com/SkipTheFish/gitops-platform/backend/internal/config"
	"github.com/SkipTheFish/gitops-platform/backend/internal/pkg"
	"github.com/SkipTheFish/gitops-platform/backend/internal/router"
)

func main() {

	// 加载配置文件
	cfg, err := config.Load()

	// 没有读取到
	if err != nil {
		// 使用 log.Fatalf() 立即退出，因为没有配置无法继续
		log.Fatalf("load config failed: %v", err)
	}

	log.Println("Config loaded")

	// 初始化数据库连接
	db, err := pkg.NewPostgres(cfg.PostgresDSN()) // 使用pkg.NewPostgres 解析 DSN 并创建 GORM 连接池

	log.Println("Postgres connected")

	// 连接失败
	if err != nil {
		log.Fatalf("connect postgres failed: %v", err)
	}

	// 初始化 Redis 客户端
	// 参数来自配置：地址、密码、DB 编号（字符串形式）
	rdb, err := pkg.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {

		// 连接失败
		log.Fatalf("connect redis failed: %v", err)
	}

	log.Println("Redis connected")

	// 创建并配置 HTTP 路由
	// 将已初始化的 db 和 rdb 作为依赖注入给 router
	// router.New() 内部会绑定 /api/health 等路径到对应的 handler
	r := router.New(db, rdb)
	log.Println("Router created")
	// 启动 HTTP 服务
	// 从配置中读取 HTTP 端口（如 "8080"），拼接成 ":8080"
	addr := ":" + cfg.HTTPPort

	// 输出启动日志
	log.Printf("server starting at %s", addr)

	// 如果端口被占用、权限不足 返回 error
	if err := r.Run(addr); err != nil {

		// 启动失败
		log.Fatalf("server run failed: %v", err)
	}
}
