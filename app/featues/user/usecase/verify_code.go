package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/middlewares"
	"mgo-gin/utils/constant"
	"net/http"
)

func VerifyCode(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyCode{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRef, code, err := userEntity.GetVerificationById(userRequest.UserRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		userRef, code, err = userEntity.ActiveVerification(userRequest.UserRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		actionToken := middlewares.GenerateActionToken(userRequest.UserRefId)
		response := gin.H{
			"actionToken": actionToken,
		}
		ctx.JSON(code, response)
	}
}
