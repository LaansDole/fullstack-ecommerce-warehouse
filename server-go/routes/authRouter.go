// routes/authRouter.go
package routes

import (
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/controllers"
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/login/whadmin", controllers.LoginAdmin)

		auth.Use(middleware.Authentication())
		auth.DELETE("/logout", controllers.Logout)
	}
}
