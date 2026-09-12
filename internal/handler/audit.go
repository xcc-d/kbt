package handler

import (
	"kbt/internal/model"
	"kbt/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	svc *service.AuditService
}

func NewAuditHandler(svc *service.AuditService) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (a *AuditHandler) List(g *gin.Context) {
	page, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(g.DefaultQuery("page_size", "20"))

	q := service.AuditQuery{
		Operator:     g.Query("operator"),
		Resource:     g.Query("resource"),
		ResourceName: g.Query("resource_name"),
		Namespace:    g.Query("namespace"),
		Action:       g.Query("action"),
		Result:       g.Query("result"),
		Page:         page,
		PageSize:     pageSize,
	}

	logs, total, err := a.svc.Audit(g.Request.Context(), q)
	if err != nil {
		model.Fail(g, err)
		return
	}
	model.Success(g, gin.H{
		"total": total,
		"page":  page,
		"items": logs,
	})
}
