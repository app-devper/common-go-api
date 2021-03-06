package usecase

import (
	"devper/app/core/constant"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyUser(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.VerifyUser{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := userEntity.GetUserByUsername(userRequest.Username)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_ = userEntity.RemoveVerificationObjective(userRequest.Objective)

		ref := form.Reference{
			UserId:    user.Id,
			Type:      constant.ActionToken,
			Objective: userRequest.Objective,
			Status:    constant.INACTIVE,
		}
		userRef, err := userEntity.CreateVerification(ref)
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
		result := gin.H{
			"userRefId":      userRef.Id,
			"verifyChannels": channels,
		}
		ctx.JSON(http.StatusOK, result)
	}
}
