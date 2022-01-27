package usecase

import (
	"devper/app/core/constant"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/config"
	"devper/middlewares"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func VerifyUserCode(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyCode{}
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
		if userRef.ExpireDate.Before(time.Now()) {
			logrus.Error(err)
			err = errors.New("token invalid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userRequest.RefId != userRef.RefId || userRequest.Code != userRef.Code {
			logrus.Error(err)
			err = errors.New("wrong code")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		expirationTime := time.Now().Add(config.ActionTokenTime)
		userRef, err = userEntity.ActiveVerification(userRequest.UserRefId, expirationTime)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		actionToken := middlewares.GenerateActionToken(userRequest.UserRefId, expirationTime)
		result := gin.H{
			"actionToken": actionToken,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
