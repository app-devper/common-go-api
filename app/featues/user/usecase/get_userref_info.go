package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/repository"
	"net/http"
)

func GetUserRefInfo(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRefId := ctx.GetString("UserRefId")
		userRef, code, err := userEntity.GetVerificationById(userRefId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, code, err := userEntity.GetUserById(userRef.UserId.Hex())
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}
