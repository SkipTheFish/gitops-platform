//

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/SkipTheFish/gitops-platform/backend/internal/service"
)

// CreateDeploymentRecordRequest 是 HTTP 层的请求体。
// 注意：它和 service 的 input 不是一回事。
// handler 的 request 负责接 JSON；service input 负责承载业务数据。
type CreateDeploymentRecordRequest struct {
	AppID               int64   `json:"app_id" binding:"required"`
	EnvID               int64   `json:"env_id" binding:"required"`
	Version             string  `json:"version" binding:"required"`
	ImageTag            string  `json:"image_tag"`
	GitCommit           string  `json:"git_commit"`
	ArgoCDAppName       string  `json:"argocd_app_name"`
	SyncStatus          string  `json:"sync_status"`
	HealthStatus        string  `json:"health_status"`
	Operator            string  `json:"operator" binding:"required"`
	RollbackFromVersion *string `json:"rollback_from_version"`
}

type DeploymentRecordHandler struct {
	service *service.DeploymentRecordService
}

func NewDeploymentRecordHandler(service *service.DeploymentRecordService) *DeploymentRecordHandler {
	return &DeploymentRecordHandler{service: service}
}

// CreateDeploymentRecord 创建一条部署记录
func (h *DeploymentRecordHandler) CreateDeploymentRecord(c *gin.Context) {
	var req CreateDeploymentRecordRequest

	// 把前端输入读取到请求体中
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// 把数据存到service的input当中
	record, err := h.service.Create(c.Request.Context(), service.CreateDeploymentRecordInput{
		AppID:               req.AppID,
		EnvID:               req.EnvID,
		Version:             req.Version,
		ImageTag:            req.ImageTag,
		GitCommit:           req.GitCommit,
		ArgoCDAppName:       req.ArgoCDAppName,
		SyncStatus:          req.SyncStatus,
		HealthStatus:        req.HealthStatus,
		Operator:            req.Operator,
		RollbackFromVersion: req.RollbackFromVersion,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create deployment record failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "deployment record created successfully",
		"data":    record,
	})
}

// GetDeploymentRecord 查询单条部署记录
func (h *DeploymentRecordHandler) GetDeploymentRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid deployment id"})
		return
	}

	record, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "deployment record not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": record})
}

// ListDeploymentsByApp 查询某个 app 的部署历史
func (h *DeploymentRecordHandler) ListDeploymentsByApp(c *gin.Context) {
	appID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid app id"})
		return
	}

	list, err := h.service.ListByAppID(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list deployments by app failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ListDeploymentsByEnv 查询某个 environment 的部署历史
func (h *DeploymentRecordHandler) ListDeploymentsByEnv(c *gin.Context) {
	envID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid environment id"})
		return
	}

	list, err := h.service.ListByEnvID(c.Request.Context(), envID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list deployments by environment failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}
