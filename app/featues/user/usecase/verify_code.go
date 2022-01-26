package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/middlewares"
	"devper/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyCode(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyCode{}
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
		if userRequest.RefId != userRef.RefId || userRequest.Code != userRef.Code {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong code"})
			return
		}
		userRef, err = userEntity.ActiveVerification(userRequest.UserRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		actionToken := middlewares.GenerateActionToken(userRequest.UserRefId)
		result := gin.H{
			"actionToken": actionToken,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
