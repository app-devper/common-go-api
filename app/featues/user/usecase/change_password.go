package usecase

import (
	"devper/app/core/constant"
	"devper/app/core/utils"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func ChangePassword(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		user, err := userEntity.GetUserByRefId(userRefId, constant.AccessApi)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		userRequest := form.ChangePassword{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if (user == nil) || utils.ComparePasswordAndHashedPassword(userRequest.OldPassword, user.Password) != nil {
			err = errors.New("wrong password")
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := userEntity.ChangePassword(user.Id.Hex(), userRequest)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
