package router

import (
	"github.com/gin-gonic/gin"

	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/handler"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/middleware"
)

func SetupRouter(orderHandler *handler.OrderHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")

	authenticated := api.Group("")
	authenticated.Use(middleware.AuthenticatedUserMiddleware())
	{
		orders := authenticated.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
		}

		authenticated.GET(
			"/user/:user_id/orders",
			orderHandler.GetUserOrders,
		)
	}

	return router
}
