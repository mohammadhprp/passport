package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter assembles application routes and injects dependencies.
func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "up"})
	})

	return router
}
