package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func VerifyUserChannel(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyChannel{}
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
		expirationTime := time.Now().Add(5 * time.Minute)
		result, err := userEntity.UpdateVerification(userRequest, expirationTime)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
