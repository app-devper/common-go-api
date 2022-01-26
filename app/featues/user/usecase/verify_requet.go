package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func VerifyRequest(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyRequest{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRef, err := userEntity.GetVerificationById(userRequest.UserRefId)
		if userRef == nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ref invalid"})
			return
		}
		if userRef.Status == constant.ACTIVE {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ref is active"})
			return
		}
		if userRef.ExpireDate.Before(time.Now()) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "token invalid"})
			return
		}
		result, err := userEntity.UpdateVerification(userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
