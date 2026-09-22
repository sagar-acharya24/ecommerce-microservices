package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/auth"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/handler"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/middleware"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	jwtManager *auth.JWTManager,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			// Public endpoints.
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)

			// Protected endpoints.
			protected := users.Group("")
			protected.Use(middleware.AuthMiddleware(jwtManager))
			{
				protected.GET("/:id", userHandler.GetUser)
				protected.PUT("/:id", userHandler.UpdateUser)
			}
		}
	}

	return router
}
