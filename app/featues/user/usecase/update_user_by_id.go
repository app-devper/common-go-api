package usecase

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/featues/user/form"
	"mgo-gin/app/featues/user/repository"
	"net/http"
)

func UpdateUserById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString("UserId")
		_, code, err := userEntity.GetUserById(userId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		id := ctx.Param("id")
		userRequest := form.UpdateUser{}
		err = ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		user, code, err := userEntity.UpdateUserById(id, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, user)
	}
}
