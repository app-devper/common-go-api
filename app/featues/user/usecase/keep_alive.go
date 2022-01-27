package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/middlewares"
	"devper/utils/constant"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func KeepAlive(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		user, err := userEntity.GetUserByRefId(userRefId, constant.AccessApi)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
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
