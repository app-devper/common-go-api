package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/repository"
)

func GetUserById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		user, code, err := userEntity.GetUserById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}
