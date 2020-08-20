package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//{
//	code: 20001,
//	data: xxx,
//	msg: xx
//}

func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, message string) {
	ctx.JSON(httpStatus, gin.H{"code": code, "data": data, "message": message})
}

func Success(ctx *gin.Context, data gin.H, message string) {
	Response(ctx, http.StatusOK, 200, data, message)
}

func Fail(ctx *gin.Context, message string, data gin.H) {
	Response(ctx, http.StatusOK, 400, data, message)
}
