package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorSender(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func NotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Not Found",
	})
}

func ResponseOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
