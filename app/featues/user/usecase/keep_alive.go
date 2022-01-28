package usecase

import (
	"devper/app/core/constant"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/config"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func KeepAlive(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		userRefId := ctx.GetString("UserRefId")
		user, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ref := form.Reference{
			UserId:      user.Id,
			Type:        constant.AccessToken,
			Objective:   constant.AccessApi,
			Channel:     "ACCESS_TOKEN",
			ChannelInfo: user.Username,
			ExpireDate:  time.Now().Add(config.AccessTokenTime),
			Status:      constant.ACTIVE,
		}
		userRef, err := userEntity.CreateVerification(ref)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, _ = userEntity.RevokeVerification(userRefId)

		token := middlewares.GenerateJwtToken(userRef.Id.Hex(), user.Role, userRef.ExpireDate)
		result := gin.H{
			"accessToken": token,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
