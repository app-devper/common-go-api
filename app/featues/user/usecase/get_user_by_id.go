package usecase

import (
	"devper/app/featues/user/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		result, err := userEntity.GetUserById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
