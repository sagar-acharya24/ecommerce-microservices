package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/api-gateway/internal/middleware"
	"github.com/sagar-acharya24/ecommerce-microservices/api-gateway/internal/proxy"
)

func SetupRouter(
	userProxy *proxy.ServiceProxy,
	productProxy *proxy.ServiceProxy,
	orderProxy *proxy.ServiceProxy,
	jwtSecret string,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	// Health check.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")

	// --------------------------------------------------
	// User Service
	// --------------------------------------------------

	// Public user endpoints.
	api.POST("/users/register", gin.WrapH(userProxy))
	api.POST("/users/login", gin.WrapH(userProxy))

	// Protected user endpoints.
	users := api.Group("/users")
	users.Use(authMiddleware)

	users.GET("/:id", gin.WrapH(userProxy))
	users.PUT("/:id", gin.WrapH(userProxy))

	// --------------------------------------------------
	// Product Service
	// --------------------------------------------------

	// Public product endpoints.
	api.GET("/products", gin.WrapH(productProxy))
	api.GET("/products/:id", gin.WrapH(productProxy))

	// Protected product endpoints.
	products := api.Group("/products")
	products.Use(authMiddleware)

	products.POST("", gin.WrapH(productProxy))
	products.PUT("/:id", gin.WrapH(productProxy))
	products.DELETE("/:id", gin.WrapH(productProxy))

	// --------------------------------------------------
	// Order Service
	// --------------------------------------------------

	// All order endpoints require authentication.
	orders := api.Group("")
	orders.Use(authMiddleware)

	orders.POST("/orders", gin.WrapH(orderProxy))
	orders.GET("/orders/:id", gin.WrapH(orderProxy))
	orders.PUT("/orders/:id/cancel", gin.WrapH(orderProxy))
	orders.GET("/user/:user_id/orders", gin.WrapH(orderProxy))

	return router
}
