// 创建一个 Redis 客户端
// 这里 Ping 了，会立即连接

package pkg

import (
	"context" // 请求传递上下文，支持超时控制和取消操作
	"strconv" // 用于字符串与基本数据类型之前的转换 (eg：str -> int)

	"github.com/redis/go-redis/v9" // Go 官方推荐的 Redis 客户端库（v9 版本）
)

func NewRedis(addr, password, dbStr string) (*redis.Client, error) {

	// 参数说明：
	//   - addr: Redis 服务器地址，格式为 "host:port"，例如 "localhost:6379"
	//   - password: Redis 认证密码（若未设置密码，可传空字符串 ""）
	//   - dbStr: 要连接的 Redis 数据库编号
	//
	// 返回值：
	//   - *redis.Client: 初始化成功且通过连通性测试的 Redis 客户端句柄
	//   - error: 初始化过程中发生的任何错误（如无效 DB 编号、连接失败、认证失败等）

	db, err := strconv.Atoi(dbStr) // 转换编号为整数

	if err != nil {
		// 转换失败，返回错误
		return nil, err
	}

	// 创建 redis 客户端实例
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,     // 服务器地址
		Password: password, // 密码
		DB:       db,       // 数据库编号
	})

	// 主动向 Redis 服务器发送 PING 命令，验证连接是否真正可用
	// 若以下任一情况发生，Ping() 将返回 error：
	//   - Redis 服务未启动或网络不通
	//   - 地址/端口错误
	//   - 密码错误
	//   - DB 编号超出范围（虽然 Redis 通常允许动态创建，但某些部署可能限制）
	if err := rdb.Ping(context.Background()).Err(); err != nil {

		// 连接验证失败，返回具体错误（如 "dial tcp: connect: connection refused"）
		// 此时不应使用该客户端，避免后续操作静默失败
		return nil, err

	}

	// 初始化成功：返回已验证的 Redis 客户端
	// 此客户端内部包含连接池，可安全地在多个 goroutine 中并发使用
	return rdb, nil
}
