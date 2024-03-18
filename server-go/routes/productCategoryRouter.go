// routes/productRouter.go
package routes

import (
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/controllers"
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/middleware"
	"github.com/gin-gonic/gin"
)

func ProductCategoryRoutes(router *gin.Engine) {
	category := router.Group("/api/product-category")
	{
		category.Use(middleware.Authentication())
		category.GET("/", controllers.GetAllProductCategories)
		category.GET("/:category_name", controllers.GetProductCategoryByName)

		category.POST("/create", middleware.CheckAdmin(), controllers.CreateProductCategory)
		category.PUT("/update/:category_name", middleware.CheckAdmin(), controllers.UpdateProduct)
		category.DELETE("/delete/:id", middleware.CheckAdmin(), controllers.DeleteProduct)
	}
}
