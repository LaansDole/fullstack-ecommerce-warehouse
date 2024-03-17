// routes/productRouter.go
package routes

import (
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/controllers"
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/middleware"
	"github.com/gin-gonic/gin"
)

func InboundOrderRoutes(router *gin.Engine) {
	inboundOrder := router.Group("/api/inbound-order")
	{
		inboundOrder.Use(middleware.Authentication(), middleware.CheckSeller())
		inboundOrder.GET("/", controllers.GetAllInboundOrders)
		inboundOrder.GET("/:id", controllers.GetInboundOrderById)
		inboundOrder.POST("/create", controllers.CreateInboundOrder)
		inboundOrder.PUT("/update/:id", controllers.UpdateInboundOrder)
		inboundOrder.DELETE("/delete/:id", controllers.DeleteInboundOrder)
		inboundOrder.POST("/fulfill/:id", controllers.FulfillInboundOrder)
	}
}
