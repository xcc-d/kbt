package handler

import (
	"io"
	"kbt/internal/model"
	"kbt/internal/service"
	apperr "kbt/pkg/errors"

	"github.com/gin-gonic/gin"
)

type DeploymentHandler struct {
	svc *service.DeploymentService
}

func NewDeploymentHandler(svc *service.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{svc: svc}
}

func (d *DeploymentHandler) List(g *gin.Context) {
	namespace := g.Param("namespace")
	list, err := d.svc.DeploymentList(g.Request.Context(), namespace)
	if err != nil {
		model.Fail(g, err)
		return
	}
	model.Success(g, list)
}

func (d *DeploymentHandler) Get(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("deployment")
	result, err := d.svc.DeploymentGet(g.Request.Context(), namespace, name)
	if err != nil {
		model.Fail(g, err)
		return
	}
	model.Success(g, result)
}

func (d *DeploymentHandler) Create(g *gin.Context) {
	reqFile, err := g.FormFile("yamlFile")
	if err != nil {
		model.Fail(g, apperr.NewBadRequest("file is required (multipart field name: yamlFile)"))
		return
	}

	file, err := reqFile.Open()
	if err != nil {
		model.Fail(g, apperr.NewFailed("failed to open uploaded file: "+err.Error()))
		return
	}
	defer file.Close()

	allFile, err := io.ReadAll(file)
	if err != nil {
		model.Fail(g, apperr.NewFailed("failed to read uploaded file: "+err.Error()))
		return
	}

	err = d.svc.DeploymentCreate(g.Request.Context(), allFile)
	if err != nil {
		model.Fail(g, err)
		return
	}
	model.Success(g, nil)
}

func (d *DeploymentHandler) Delete(g *gin.Context) {
	namespace := g.Param("namespace")
	name := g.Param("deployment")
	err := d.svc.DeploymentDelete(g.Request.Context(), namespace, name)
	if err != nil {
		model.Fail(g, err)
		return
	}
	model.Success(g, nil)
}
