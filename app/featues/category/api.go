package category

import (
	"devper/app/featues/category/form"
	"devper/app/featues/category/repository"
	"devper/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApplyCategoryAPI(app *gin.RouterGroup, categoryEntity repository.ICategory) {
	productRoute := app.Group("category")

	productRoute.GET("",
		getCategoryAll(categoryEntity),
	)

	productRoute.POST("",
		middlewares.RequireAuthenticated(),
		createCategory(categoryEntity),
	)

	productRoute.GET("/:categoryId",
		getCategoryById(categoryEntity),
	)

	productRoute.PUT("/:categoryId",
		middlewares.RequireAuthenticated(),
		updateCategoryById(categoryEntity),
	)

	productRoute.DELETE("/:categoryId",
		middlewares.RequireAuthenticated(),
		deleteCategoryById(categoryEntity),
	)

	productRoute.PATCH("/:categoryId/default",
		middlewares.RequireAuthenticated(),
		updateDefaultCategoryById(categoryEntity),
	)
}

func deleteCategoryById(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		categoryId := ctx.Param("categoryId")
		result, code, err := entity.RemoveCategoryById(categoryId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func updateCategoryById(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		categoryId := ctx.Param("categoryId")
		request := form.Category{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, code, err := entity.UpdateCategoryById(categoryId, request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func updateDefaultCategoryById(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		categoryId := ctx.Param("categoryId")
		request := form.Category{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, code, err := entity.UpdateDefaultCategoryById(categoryId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func getCategoryById(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		categoryId := ctx.Param("categoryId")
		result, code, err := entity.GetCategoryById(categoryId)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func createCategory(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := form.Category{}
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, code, err := entity.CreateCategory(request)
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}

func getCategoryAll(entity repository.ICategory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, code, err := entity.GetCategoryAll()
		if err != nil {
			ctx.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(code, result)
	}
}
