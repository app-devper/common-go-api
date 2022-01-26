package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"net/http"
)

func VerifyUser(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyUser{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, _, err := userEntity.GetUserByUsername(userRequest.Username)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userRef, code, err := userEntity.CreateVerification(user.Id.Hex(), userRequest.Objective)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		channels := []form.Channel{{
			Channel:     "MOBILE",
			ChannelInfo: user.Phone,
		}, {
			Channel:     "EMAIL",
			ChannelInfo: user.Email,
		}}
		response := gin.H{
			"userRefId":      userRef.Id,
			"verifyChannels": channels,
		}
		ctx.JSON(code, response)
	}
}
