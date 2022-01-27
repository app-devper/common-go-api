package usecase

import (
	"devper/app/core/constant"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func UpdateRoleById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		user, err := userEntity.GetUserByRefId(userRefId, constant.AccessApi)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		userId := user.Id.Hex()
		if userId == id {
			err = errors.New("can't update self user")
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't update self user"})
			return
		}
		userRequest := form.UpdateRole{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		result, err := userEntity.UpdateRoleById(id, userRequest)
		if err != nil {
			logrus.Error(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
