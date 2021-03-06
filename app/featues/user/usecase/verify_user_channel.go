package usecase

import (
	"devper/app/core/constant"
	"devper/app/core/utils"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/config"
	"errors"
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
			err = errors.New("user ref invalid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userRef.Status == constant.ACTIVE {
			err = errors.New("user ref is active")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		expirationTime := time.Now().Add(config.VerifyCodeTime)
		result, err := userEntity.UpdateVerification(userRequest, expirationTime)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, _ = utils.NotifyMassage("RefId :" + result.RefId + " Code :" + result.Code)

		ctx.JSON(http.StatusOK, result)
	}
}
