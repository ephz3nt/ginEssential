package middleware

import (
	"net/http"
	"strings"

	"github.com/ephz3nt/ginEssential/common"
	"github.com/ephz3nt/ginEssential/model"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get authorization header
		tokenString := ctx.GetHeader("Authorization")

		// validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "permission denied"})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "permission denied"})
			ctx.Abort()
			return
		}

		// validate success, get claim userid
		userID := claims.UserID
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userID)

		// 验证用户是否存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "permission denied"})
			ctx.Abort()
			return
		}

		// 用户存在 ， 将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
