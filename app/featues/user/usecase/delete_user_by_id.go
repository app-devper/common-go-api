package usecase

import (
	"devper/app/featues/user/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteUserById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		id := ctx.Param("id")
		if userId == id {
			err := errors.New("can't delete self user")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := userEntity.RemoveUserById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
