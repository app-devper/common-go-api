package usecase

import (
	"devper/app/core/constant"
	"devper/app/core/utils"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/config"
	"devper/middlewares"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func VerifyPassword(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyPassword{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId := ctx.GetString("UserId")
		user, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if (user == nil) || utils.ComparePasswordAndHashedPassword(userRequest.Password, user.Password) != nil {
			err = errors.New("wrong password")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_ = userEntity.RemoveVerificationObjective(userRequest.Objective)

		ref := form.Reference{
			UserId:      user.Id,
			Type:        constant.ActionToken,
			Objective:   userRequest.Objective,
			Channel:     "USERNAME",
			ChannelInfo: user.Username,
			ExpireDate:  time.Now().Add(config.ActionTokenTime),
			Status:      constant.ACTIVE,
		}
		userRef, err := userEntity.CreateVerification(ref)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := middlewares.GenerateActionToken(userRef.Id.Hex(), userRef.Objective, userRef.ExpireDate)
		result := gin.H{
			"actionToken": token,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
