package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/SkipTheFish/gitops-platform/backend/internal/service"
)

// CreatePipelineRunRequest 是 HTTP 层接收的“手动触发发布”请求体。
type CreatePipelineRunRequest struct {
	AppID       int64  `json:"app_id" binding:"required"`
	EnvID       int64  `json:"env_id" binding:"required"`
	GitCommit   string `json:"git_commit"`
	Branch      string `json:"branch"`
	ImageTag    string `json:"image_tag" binding:"required"`
	TriggerType string `json:"trigger_type"`
	Operator    string `json:"operator" binding:"required"`
	Version     string `json:"version" binding:"required"`
}

type PipelineRunHandler struct {
	pipelineService *service.PipelineRunService
}

func NewPipelineRunHandler(pipelineService *service.PipelineRunService) *PipelineRunHandler {
	return &PipelineRunHandler{
		pipelineService: pipelineService,
	}
}

// CreateManualRun 手动触发一次发布任务
func (h *PipelineRunHandler) CreateManualRun(c *gin.Context) {
	var req CreatePipelineRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	run, err := h.pipelineService.CreateManualRun(c.Request.Context(), service.CreatePipelineRunInput{
		AppID:       req.AppID,
		EnvID:       req.EnvID,
		GitCommit:   req.GitCommit,
		Branch:      req.Branch,
		ImageTag:    req.ImageTag,
		TriggerType: req.TriggerType,
		Operator:    req.Operator,
		Version:     req.Version,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create pipeline run failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "pipeline run created successfully",
		"data":    run,
	})
}

// GetPipelineRun 查询单条流水线任务
func (h *PipelineRunHandler) GetPipelineRun(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid pipeline id"})
		return
	}

	run, err := h.pipelineService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "pipeline run not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": run})
}

// ListPipelineRunsByApp 查询某个 app 的流水线历史
func (h *PipelineRunHandler) ListPipelineRunsByApp(c *gin.Context) {
	appID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid app id"})
		return
	}

	list, err := h.pipelineService.ListByAppID(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list pipeline runs by app failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// ListPipelineRunsByEnv 查询某个环境的流水线历史
func (h *PipelineRunHandler) ListPipelineRunsByEnv(c *gin.Context) {
	envID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid environment id"})
		return
	}

	list, err := h.pipelineService.ListByEnvID(c.Request.Context(), envID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list pipeline runs by environment failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}
