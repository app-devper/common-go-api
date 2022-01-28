package usecase

import (
	"devper/app/featues/user/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := userEntity.GetUserAll()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
