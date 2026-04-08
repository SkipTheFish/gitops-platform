package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/SkipTheFish/gitops-platform/backend/internal/service"
)

type CreateOperationAuditRequest struct {
	Operator   string `json:"operator" binding:"required"`
	ActionType string `json:"action_type" binding:"required"`
	TargetID   int64  `json:"target_id" binding:"required"`
	Detail     string `json:"detail"`
}

type OperationAuditHandler struct {
	service *service.OperationAuditService
}

func NewOperationAuditHandler(service *service.OperationAuditService) *OperationAuditHandler {
	return &OperationAuditHandler{service: service}
}

// CreateAudit 手动创建一条审计记录
// 虽然后面很多审计会由系统自动生成，但保留这个接口有助于你调试和学习。
func (h *OperationAuditHandler) CreateAudit(c *gin.Context) {
	var req CreateOperationAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	audit, err := h.service.Create(c.Request.Context(), service.CreateOperationAuditInput{
		Operator:   req.Operator,
		ActionType: req.ActionType,
		TargetID:   req.TargetID,
		Detail:     req.Detail,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create audit failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "audit created successfully",
		"data":    audit,
	})
}

// ListAuditsByTargetID 查询某个目标对象对应的审计日志
func (h *OperationAuditHandler) ListAuditsByTargetID(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid target id"})
		return
	}

	list, err := h.service.ListByTargetID(c.Request.Context(), targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list audits failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}
