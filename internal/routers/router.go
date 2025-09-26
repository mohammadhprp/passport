package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mohammadhprp/passport/internal/handlers"
	"github.com/mohammadhprp/passport/internal/repositories"
	"github.com/mohammadhprp/passport/internal/services"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "up"})
	})

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	userRoutes := router.Group("/users")
	userHandler.RegisterRoutes(userRoutes)

	return router
}
