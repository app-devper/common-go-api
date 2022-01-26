package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/repository"
	"net/http"
)

func DeleteUserById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		if userId == id {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't delete self user"})
			return
		}
		user, code, err := userEntity.RemoveUserById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}
