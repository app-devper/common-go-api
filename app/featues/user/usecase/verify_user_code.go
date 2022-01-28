package usecase

import (
	"devper/app/core/constant"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/config"
	"devper/middlewares"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func VerifyUserCode(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyCode{}
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
		if userRef.ExpireDate.Before(time.Now()) {
			err = errors.New("token invalid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if userRequest.RefId != userRef.RefId || userRequest.Code != userRef.Code {
			err = errors.New("code invalid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRef, err = userEntity.ActiveVerification(userRequest.UserRefId, time.Now().Add(config.ActionTokenTime))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		actionToken := middlewares.GenerateActionToken(userRef.Id.Hex(), userRef.Objective, userRef.ExpireDate)
		result := gin.H{
			"actionToken": actionToken,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
