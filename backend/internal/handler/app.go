package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/SkipTheFish/gitops-platform/backend/internal/service"
)

// CreateAppRequest 表示“新建应用”的请求体。
// binding:"required" 表示 Gin 在绑定 JSON 时会校验该字段必填。
// ->	当请求到达后端时，你会调用类似 c.ShouldBindJSON(&req) 的方法。Gin 会查看这个结构体标签（Tag）：
// ->	如果带这个（binding:"required"）tag ，前端只要有字段没写数值，就会拦截报错，避免后续
// ->	这种写法的核心价值在于防御式编程。它保证了进入你业务逻辑层（Service 层）的数据一定是完整、合法的，
// ->	避免了在代码里写一大堆 if name == "" { ... } 这种啰嗦的判断。
type CreateAppRequest struct {
	Name           string `json:"name" binding:"required"`
	RepoURL        string `json:"repo_url" binding:"required"`
	ConfigRepoURL  string `json:"config_repo_url" binding:"required"`
	ClusterName    string `json:"cluster_name" binding:"required"`
	Namespace      string `json:"namespace" binding:"required"`
	HelmChartPath  string `json:"helm_chart_path" binding:"required"`
	ValuesFilePath string `json:"values_file_path" binding:"required"`
}

// UpdateAppRequest 表示“修改应用”的请求体。
// 修改接口通常不要求所有字段必填，因此这里不加 required。
type UpdateAppRequest struct {
	Name           string `json:"name"`
	RepoURL        string `json:"repo_url"`
	ConfigRepoURL  string `json:"config_repo_url"`
	ClusterName    string `json:"cluster_name"`
	Namespace      string `json:"namespace"`
	HelmChartPath  string `json:"helm_chart_path"`
	ValuesFilePath string `json:"values_file_path"`
}

// 一个service结构体
type AppHandler struct {
	appService *service.AppService
}

// 构造函数
// 输入：*service.AppService（AppService 的指针）
// 输出：*AppHandler（AppHandler 的指针）
// NewAppHandler 用依赖注入的方式把 service 注入进来。
// 可以把它理解成：handler 不自己创建 service，而是外部组装好再传进来。
func NewAppHandler(appService *service.AppService) *AppHandler {
	return &AppHandler{
		appService: appService,
	}
}

// CreateApp 新增应用
// 输入：指向AppHandler的指针h、指向gin.Context的指针c
func (h *AppHandler) CreateApp(c *gin.Context) {

	// 创建一个请求体
	var req CreateAppRequest

	// ShouldBindJSON 去读取 HTTP 请求里的 Body（那一串 JSON 文本），检查请求体是否合法，不合法引出异常
	// 	内部逻辑：ShouldBindJSON 会去读取 HTTP 请求里的 Body（那一串 JSON 文本）。
	// 反射赋值：它通过 Go 的**反射（Reflection）**机制，
	// 			对比 JSON 里的 Key（比如 "name": "my-new-app"）和你的结构体标签 json:"name"。
	// 装填数据：如果对上了，它就直接把 JSON 里的值填进 req 对象的内存空间里。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	//  DTO 转 Input 结构体
	// 就是把请求体的数据读取到service.CreateAppInput的结构体当中
	input := service.CreateAppInput{
		Name:           req.Name,
		RepoURL:        req.RepoURL,
		ConfigRepoURL:  req.ConfigRepoURL,
		ClusterName:    req.ClusterName,
		Namespace:      req.Namespace,
		HelmChartPath:  req.HelmChartPath,
		ValuesFilePath: req.ValuesFilePath,
	}

	// 调用appService来创建
	app, err := h.appService.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create app failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "app created successfully",
		"data":    app,
	})
}

// ListApps 查询应用列表
func (h *AppHandler) ListApps(c *gin.Context) {
	// 调用appService中的 list 方法
	apps, err := h.appService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "list apps failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": apps,
	})
}

// GetApp 查询单个应用详情
func (h *AppHandler) GetApp(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64) // 转换成10进制，64位宽
	// 转换失败
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid app id",
		})
		return
	}

	// 调用appService中的 获取单独id 方法
	app, err := h.appService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "app not found",
			"error":   err.Error(),
		})
		return
	}

	// 将数据发送回前端
	// http.StatusOK：即 HTTP 状态码 200，告诉前端“一切正常”。
	// gin.H：这是 Gin 提供的便捷写法，
	// 本质上是一个 map[string]interface{}（类似 Python 的字典或 JS 的对象）。
	c.JSON(http.StatusOK, gin.H{
		"data": app,
		// 序列化：Gin 会自动把 app（这个 Go 结构体）转换成 JSON 字符串，
		// 并设置 HTTP Header 为 Content-Type: application/json。
	})
}

// UpdateApp 修改应用
func (h *AppHandler) UpdateApp(c *gin.Context) {
	// 依旧转换id
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	// 检查id是否合法
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid app id",
		})
		return
	}

	// 从前段读取数据到请求体
	var req UpdateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// 赋值到 input
	input := service.UpdateAppInput{
		ID:             id,
		Name:           req.Name,
		RepoURL:        req.RepoURL,
		ConfigRepoURL:  req.ConfigRepoURL,
		ClusterName:    req.ClusterName,
		Namespace:      req.Namespace,
		HelmChartPath:  req.HelmChartPath,
		ValuesFilePath: req.ValuesFilePath,
	}

	// 调用appService，进行唯一性检验
	app, err := h.appService.Update(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update app failed",
			"error":   err.Error(),
		})
		return
	}

	// 更新完成，将数据发送回前端
	c.JSON(http.StatusOK, gin.H{
		"message": "app updated successfully",
		"data":    app,
	})
}
