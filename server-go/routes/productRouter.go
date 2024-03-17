// routes/productRouter.go
package routes

import (
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/controllers"
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/middleware"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	product := router.Group("/api/product")
	{
		product.Use(middleware.Authentication())
		product.GET("/", controllers.GetAllProducts)
		product.GET("/asc", controllers.GetAllProductsASC)
		product.GET("/dsc", controllers.GetAllProductsDSC)
		product.GET("/:id", controllers.GetProductById)
		product.GET("/title/:title", controllers.GetProductByTitle)

		product.POST("/create", middleware.CheckSeller(), middleware.FileHandling(), controllers.CreateProduct)
	}
}
