package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/SkipTheFish/gitops-platform/backend/internal/service"
)

type CreateEnvironmentRequest struct {
	EnvName         string `json:"env_name" binding:"required"`
	ClusterName     string `json:"cluster_name" binding:"required"`
	Namespace       string `json:"namespace" binding:"required"`
	AutoSyncEnabled bool   `json:"auto_sync_enabled"`
	ValuesFilePath  string `json:"values_file_path"`
	ArgoCDAppName   string `json:"argocd_app_name"`
}

type UpdateEnvironmentRequest struct {
	EnvName         string `json:"env_name"`
	ClusterName     string `json:"cluster_name"`
	Namespace       string `json:"namespace"`
	AutoSyncEnabled *bool  `json:"auto_sync_enabled"`
	ValuesFilePath  string `json:"values_file_path"`
	ArgoCDAppName   string `json:"argocd_app_name"`
}

type EnvironmentHandler struct {
	envService *service.EnvironmentService
}

func NewEnvironmentHandler(envService *service.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{
		envService: envService,
	}
}

// CreateEnvironment 为指定 app 创建环境
func (h *EnvironmentHandler) CreateEnvironment(c *gin.Context) {
	appID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid app id",
		})
		return
	}

	var req CreateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	input := service.CreateEnvironmentInput{
		AppID:           appID,
		EnvName:         req.EnvName,
		ClusterName:     req.ClusterName,
		Namespace:       req.Namespace,
		AutoSyncEnabled: req.AutoSyncEnabled,
		ValuesFilePath:  req.ValuesFilePath,
		ArgoCDAppName:   req.ArgoCDAppName,
	}

	env, err := h.envService.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create environment failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "environment created successfully",
		"data":    env,
	})
}

// ListEnvironments 查询某个 app 下的环境列表
func (h *EnvironmentHandler) ListEnvironments(c *gin.Context) {
	appID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid app id",
		})
		return
	}

	list, err := h.envService.ListByAppID(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "list environments failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

// GetEnvironment 查询单个环境
func (h *EnvironmentHandler) GetEnvironment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid environment id",
		})
		return
	}

	env, err := h.envService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "environment not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": env,
	})
}

// UpdateEnvironment 修改环境
func (h *EnvironmentHandler) UpdateEnvironment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid environment id",
		})
		return
	}

	var req UpdateEnvironmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	input := service.UpdateEnvironmentInput{
		ID:              id,
		EnvName:         req.EnvName,
		ClusterName:     req.ClusterName,
		Namespace:       req.Namespace,
		AutoSyncEnabled: req.AutoSyncEnabled,
		ValuesFilePath:  req.ValuesFilePath,
		ArgoCDAppName:   req.ArgoCDAppName,
	}

	env, err := h.envService.Update(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update environment failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "environment updated successfully",
		"data":    env,
	})
}
