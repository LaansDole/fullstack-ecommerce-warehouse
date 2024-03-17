// routes/productRouter.go
package routes

import (
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/controllers"
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/middleware"
	"github.com/gin-gonic/gin"
)

func StockpileRoutes(router *gin.Engine) {
	stockpile := router.Group("/api/stock")
	{
		stockpile.Use(middleware.Authentication())
		stockpile.GET("/", controllers.GetAllStockpile)
	}
}
