package usecase

import (
	"devper/app/featues/user/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetVerifyInfo(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		userRef, _ := userEntity.GetVerificationById(userRefId)
		if userRef == nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ref invalid"})
			return
		}
		ctx.JSON(http.StatusOK, userRef)
	}
}
