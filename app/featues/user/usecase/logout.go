package usecase

import (
	"devper/app/featues/user/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		_, _ = userEntity.RevokeVerification(userRefId)
		result := gin.H{
			"message": "success",
		}
		ctx.JSON(http.StatusOK, result)
	}
}
