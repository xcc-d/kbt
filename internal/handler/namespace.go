package handler

import (
	"kbt/internal/model"
	"kbt/internal/service"
	apperr "kbt/pkg/errors"

	"github.com/gin-gonic/gin"
)

type NamespaceHandler struct {
	svc *service.NamespaceService
}

func NewNamespaceHandler(svc *service.NamespaceService) *NamespaceHandler {
	return &NamespaceHandler{svc: svc}
}

func (h *NamespaceHandler) List(g *gin.Context) {
	list, err := h.svc.NamespaceList(g.Request.Context())
	if err != nil {
		Fail(g, err)
		return
	}
	Success(g, list)
}

func (h *NamespaceHandler) Get(g *gin.Context) {
	name := g.Param("namespace")
	ns, err := h.svc.NamespaceGet(g.Request.Context(), name)
	if err != nil {
		Fail(g, err)
		return
	}
	Success(g, ns)
}

func (h *NamespaceHandler) Create(g *gin.Context) {
	var req model.Namespace
	if err := g.ShouldBindJSON(&req); err != nil {
		Fail(g, apperr.NewBadRequest("name is required"))
		return
	}
	create, err := h.svc.NamespaceCreate(g.Request.Context(), &req)
	if err != nil {
		Fail(g, err)
		return
	}
	Success(g, create)
}

func (h *NamespaceHandler) Delete(g *gin.Context) {
	err := h.svc.NamespaceDelete(g.Request.Context(), g.Param("namespace"))
	if err != nil {
		Fail(g, err)
		return
	}
	Success(g, nil)
}
