package usecase

import (
	"devper/app/core/constant"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/config"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func KeepAlive(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		user, err := userEntity.GetUserByRefId(userRefId, constant.AccessApi)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		expirationTime := time.Now().Add(config.AccessTokenTime)
		ref := form.Reference{
			UserId:      user.Id,
			Objective:   constant.AccessApi,
			Channel:     "USERNAME",
			ChannelInfo: user.Username,
			ExpireDate:  expirationTime,
			Status:      constant.ACTIVE,
		}
		userRef, err := userEntity.CreateVerification(ref)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token := middlewares.GenerateJwtToken(userRef.Id.Hex(), user.Role, expirationTime)
		result := gin.H{
			"accessToken": token,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
