package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateRoleById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.UpdateRole{}
		err := ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId := ctx.GetString("UserId")
		id := ctx.Param("id")
		if userId == id {
			err := errors.New("can't update self user")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userRequest.UpdatedBy = userId
		result, err := userEntity.UpdateRoleById(id, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
