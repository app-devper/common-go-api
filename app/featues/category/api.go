package category

import (
	"devper/app/featues/category/repository"
	"devper/app/featues/category/usecase"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
)

func ApplyCategoryAPI(app *gin.RouterGroup, categoryEntity repository.ICategory) {
	productRoute := app.Group("category")

	productRoute.GET("",
		middlewares.RequireAuthenticated(),
		usecase.GetCategories(categoryEntity),
	)

	productRoute.POST("",
		middlewares.RequireAuthenticated(),
		usecase.CreateCategory(categoryEntity),
	)

	productRoute.GET("/:categoryId",
		middlewares.RequireAuthenticated(),
		usecase.GetCategoryById(categoryEntity),
	)

	productRoute.PUT("/:categoryId",
		middlewares.RequireAuthenticated(),
		usecase.UpdateCategoryById(categoryEntity),
	)

	productRoute.DELETE("/:categoryId",
		middlewares.RequireAuthenticated(),
		usecase.DeleteCategoryById(categoryEntity),
	)

	productRoute.PATCH("/:categoryId/default",
		middlewares.RequireAuthenticated(),
		usecase.UpdateDefaultCategoryById(categoryEntity),
	)
}
