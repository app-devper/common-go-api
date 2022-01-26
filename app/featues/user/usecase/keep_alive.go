package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/middlewares"
	"mgo-gin/utils/constant"
	"net/http"
)

func KeepAlive(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		user, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		if user.Status != constant.ACTIVE {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not active"})
			return
		}
		token := middlewares.GenerateJwtToken(*user)
		response := gin.H{
			"accessToken": token,
		}
		ctx.JSON(code, response)
	}
}
