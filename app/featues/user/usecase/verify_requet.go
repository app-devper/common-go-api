package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"mgo-gin/utils/constant"
	"net/http"
)

func VerifyRequest(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyRequest{}
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
		userRef, code, err = userEntity.UpdateVerification(userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, userRef)
	}
}
