package category

import (
	"devper/app/featues/category/repository"
	"devper/app/featues/category/usecase"
	repository2 "devper/app/featues/user/repository"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
)

func ApplyCategoryAPI(
	app *gin.RouterGroup,
	categoryEntity repository.ICategory,
	userEntity repository2.IUser,
) {
	productRoute := app.Group("category")

	productRoute.GET("",
		middlewares.RequireAuthenticated(userEntity),
		usecase.GetCategories(categoryEntity),
	)

	productRoute.POST("",
		middlewares.RequireAuthenticated(userEntity),
		usecase.CreateCategory(categoryEntity),
	)

	productRoute.GET("/:categoryId",
		middlewares.RequireAuthenticated(userEntity),
		usecase.GetCategoryById(categoryEntity),
	)

	productRoute.PUT("/:categoryId",
		middlewares.RequireAuthenticated(userEntity),
		usecase.UpdateCategoryById(categoryEntity),
	)

	productRoute.DELETE("/:categoryId",
		middlewares.RequireAuthenticated(userEntity),
		usecase.DeleteCategoryById(categoryEntity),
	)

	productRoute.PATCH("/:categoryId/default",
		middlewares.RequireAuthenticated(userEntity),
		usecase.UpdateDefaultCategoryById(categoryEntity),
	)
}
