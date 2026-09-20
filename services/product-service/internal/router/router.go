package router

import (
	"github.com/gin-gonic/gin"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/handler"
)

func SetupRouter(productHandler *handler.ProductHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// The service is currently accessed directly and is not behind
	// a reverse proxy, so do not trust any proxy headers.
	if err := router.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	api := router.Group("/api/v1")
	{
		products := api.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	return router
}
