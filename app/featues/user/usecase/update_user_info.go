package usecase

import (
	"devper/app/core/constant"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func UpdateUserInfo(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		user, err := userEntity.GetUserByRefId(userRefId, constant.AccessApi)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		userId := user.Id.Hex()
		userRequest := form.UpdateUser{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		result, err := userEntity.UpdateUserById(userId, userRequest)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
