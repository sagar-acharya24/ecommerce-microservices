package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/services/order-service/internal/handler"
)

func SetupRouter(orderHandler *handler.OrderHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if err := router.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	api := router.Group("/api/v1")
	{
		orders := api.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
		}

		api.GET(
			"/user/:user_id/orders",
			orderHandler.GetUserOrders,
		)
	}

	return router
}
