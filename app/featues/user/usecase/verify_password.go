package usecase

import (
	"devper/app/core/bcrypt"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/middlewares"
	"devper/utils/constant"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func VerifyPassword(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		user, err := userEntity.GetUserByRefId(userRefId, constant.AccessApi)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.VerifyPassword{}
		if err = ctx.ShouldBind(&userRequest); err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if (user == nil) || bcrypt.ComparePasswordAndHashedPassword(userRequest.Password, user.Password) != nil {
			logrus.Error(err)
			err = errors.New("wrong password")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		_ = userEntity.RemoveVerificationObjective(userRequest.Objective)
		expirationTime := time.Now().Add(3 * time.Minute)
		ref := form.Reference{
			UserId:      user.Id,
			Objective:   userRequest.Objective,
			Channel:     "ACCESS_TOKEN",
			ChannelInfo: userRefId,
			ExpireDate:  expirationTime,
			Status:      constant.ACTIVE,
		}
		userRef, err := userEntity.CreateVerification(ref)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token := middlewares.GenerateActionToken(userRef.Id.Hex(), expirationTime)
		result := gin.H{
			"actionToken": token,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
