package model

import (
	"net/http"

	apperr "kbt/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TraceId string      `json:"trace_id"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    apperr.Success,
		Message: "success",
		Data:    data,
	})
}

func Fail(c *gin.Context, err error) {
	appErr := apperr.FromError(err)
	c.JSON(appErr.Code, Response{
		Code:    appErr.Code,
		Message: appErr.Message,
		Data:    nil,
	})
}

func FailAbort(c *gin.Context, err error) {
	appErr := apperr.FromError(err)
	c.AbortWithStatusJSON(appErr.Code, Response{
		Code:    appErr.Code,
		Message: appErr.Message,
		Data:    nil,
	})
}
