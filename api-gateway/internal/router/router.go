package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/api-gateway/internal/proxy"
)

func SetupRouter(
	userProxy *proxy.ServiceProxy,
	productProxy *proxy.ServiceProxy,
	orderProxy *proxy.ServiceProxy,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		// Order Service route for user's orders.
		api.GET("/user/:user_id/orders", gin.WrapH(orderProxy))

		api.Any("/users", gin.WrapH(userProxy))
		api.Any("/users/*path", gin.WrapH(userProxy))

		api.Any("/products", gin.WrapH(productProxy))
		api.Any("/products/*path", gin.WrapH(productProxy))

		api.Any("/orders", gin.WrapH(orderProxy))
		api.Any("/orders/*path", gin.WrapH(orderProxy))
	}

	return router
}
