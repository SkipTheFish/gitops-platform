// 实现服务健康检查端点（/healthz 或 /health）。

package handler

import (
	"context"  // 用于传递取消信号和超时控制（尤其在 Redis Ping 中）
	"net/http" // 提供 HTTP 状态码常量（如 http.StatusOK）
	"time"     // 用于返回当前时间戳，便于监控系统判断时钟同步状态

	"github.com/gin-gonic/gin"     // Gin Web 框架，用于处理 HTTP 请求
	"github.com/redis/go-redis/v9" // Redis 客户端库
	"gorm.io/gorm"                 // GORM ORM 库，用于数据库操作
)

// 定义结构体
type HealthHandler struct {
	DB    *gorm.DB      // 已经初始化的 PostgreSQL
	Redis *redis.Client // 已经验证的 redis 客户端
}

// HealthHandler 是健康检查处理器，依赖数据库和 Redis 客户端。
// 采用“依赖注入”设计：外部传入 DB 和 Redis 实例，而非在内部创建。
// 这样做有三大好处：
//   1. 解耦：handler 不关心如何初始化 DB/Redis
//   2. 可测试：单元测试时可传入 mock 对象
//   3. 复用：使用已在 main.go 中初始化好的全局单例

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	// 是 HealthHandler 的构造函数
	// 参数：
	// 		HeathHandler 结构体
	// 返回
	// 		依旧结构体

	return &HealthHandler{
		DB:    db,
		Redis: redis,
	}
}

// 方法
func (h *HealthHandler) Health(c *gin.Context) {
	// 参数：
	// 		HealthHandler
	// 返回：

	type item struct {
		Status string `json:"status"`          // 组件状态："ok" 或 "down"
		Error  string `json:"error,omitempty"` // 可选错误信息（仅当状态异常时存在）

		// 内部结构体 item 用于统一表示每个组件的健康状态
		// 使用小写字段名（非导出）因为只在本函数内使用
		// json tag 控制序列化行为：
		//   - status 总是返回
		//   - error 仅在存在时返回（omitempty）
	}

	resp := gin.H{
		// 初始化默认响应：假设所有组件都正常
		// time 字段提供当前服务器时间，可用于检测时钟漂移或请求延迟

		"app":   item{Status: "ok"}, // 假设正常
		"db":    item{Status: "ok"},
		"redis": item{Status: "ok"},
		"time":  time.Now().Format(time.RFC3339), // ISO 8601 时间格式，如 "2026-04-03T12:08:00Z"
	}

	// 🔍 检查 PostgreSQL 数据库连通性
	//	注意：GORM 的 *gorm.DB 不直接提供 Ping() 方法，这里先获取底层 *sql.DB
	sqlDB, err := h.DB.DB() // 获取原生 database/sql 的 *sql.DB 对象
	if err != nil || sqlDB.Ping() != nil {

		// 不通
		resp["db"] = item{Status: "down", Error: "postgres unavailable"}
	}

	// 🔍 检查 Redis 连通性
	// go-redis/v9 的 Ping() 方法需要 context（支持超时/取消）
	// 使用 context.Background() 表示无超时限制（生产环境建议加超时！）
	if err := h.Redis.Ping(context.Background()).Err(); err != nil {

		// Redis 连接失败：可能是服务未启动、密码错误、网络问题等
		// 这里选择暴露原始错误信息（err.Error()），便于调试
		resp["redis"] = item{Status: "down", Error: err.Error()}
	}

	// 返回 JSON 响应，HTTP 状态码始终为 200
	// 即使某些组件 "down"，也返回 200 —— 这是“Liveness Probe” vs “Readiness Probe”的关键区别
	//   - Liveness（存活探针）：进程是否活着 → 失败则重启 Pod
	//   - Readiness（就绪探针）：是否准备好接收流量 → 失败则从负载均衡摘除
	// 此处实现的是 **Readiness Check**，所以即使 DB 挂了也不该返回 5xx
	c.JSON(http.StatusOK, resp)
}
