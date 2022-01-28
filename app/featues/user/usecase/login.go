package usecase

import (
	"devper/app/core/constant"
	"devper/app/core/utils"
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"devper/middlewares"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Login(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.Login{}
		if err := ctx.ShouldBind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := userEntity.GetUserByUsername(userRequest.Username)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if (user == nil) || utils.ComparePasswordAndHashedPassword(userRequest.Password, user.Password) != nil {
			err = errors.New("wrong username or password")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		expirationTime := time.Now().Add(24 * time.Hour)
		ref := form.Reference{
			UserId:      user.Id,
			Type:        constant.AccessToken,
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
