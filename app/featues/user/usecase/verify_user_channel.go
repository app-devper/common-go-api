package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/utils/constant"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func VerifyUserChannel(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyChannel{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRef, err := userEntity.GetVerificationById(userRequest.UserRefId)
		if userRef == nil {
			logrus.Error(err)
			err = errors.New("user ref invalid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userRef.Status == constant.ACTIVE {
			logrus.Error(err)
			err = errors.New("user ref is active")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		expirationTime := time.Now().Add(5 * time.Minute)
		result, err := userEntity.UpdateVerification(userRequest, expirationTime)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
