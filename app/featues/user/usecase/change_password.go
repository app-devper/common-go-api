package usecase

import (
	"devper/app/core/utils"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ChangePassword(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.ChangePassword{}
		err := ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId := ctx.GetString("UserId")
		user, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if (user == nil) || utils.ComparePasswordAndHashedPassword(userRequest.OldPassword, user.Password) != nil {
			err = errors.New("wrong password")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := userEntity.ChangePassword(user.Id.Hex(), userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
