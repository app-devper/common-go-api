package usecase

import (
	"devper/app/featues/user/form"
	"devper/app/featues/user/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateUserById(userEntity repository.IUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRequest := form.UpdateUser{}
		err := ctx.ShouldBind(&userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId := ctx.GetString("UserId")
		id := ctx.Param("id")
		userRequest.UpdatedBy = userId
		result, err := userEntity.UpdateUserById(id, userRequest)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
