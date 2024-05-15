package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
}

func (a *AuthHandler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hello world")
	return
}
