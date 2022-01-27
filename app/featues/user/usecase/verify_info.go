package usecase

import (
	"devper/app/featues/user/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetVerifyInfo(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("verifyId")
		userRef, err := userEntity.GetVerificationById(userRefId)
		if userRef == nil {
			logrus.Error(err)
			err = errors.New("user ref invalid")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, userRef)
	}
}
