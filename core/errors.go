package corev2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func StatusBadRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message": err,
		"code":    "BAD_REQUEST",
	})
}

func StatusInternalServerError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": err,
		"code":    "INTERNAL_SERVER_ERROR",
	})
}

func StatusTooManyRequests(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusTooManyRequests, gin.H{
		"message": err,
		"code":    "TOO_MANY_REQUESTS",
	})
}
