package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohammadhprp/passport/internal/handlers"
	"github.com/mohammadhprp/passport/internal/repositories"
	"github.com/mohammadhprp/passport/internal/services"
	"gorm.io/gorm"
)

// NewRouter assembles application routes and injects dependencies.
func NewRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "up"})
	})

	tenantRepo := repositories.NewTenantRepository(db)
	tenantService := services.NewTenantService(tenantRepo)
	tenantHandler := handlers.NewTenantHandler(tenantService)

	tenants := router.Group("/tenants")
	{
		tenants.GET("", tenantHandler.ListTenants)
		tenants.POST("", tenantHandler.CreateTenant)
		tenants.GET("/:id", tenantHandler.GetTenant)
		tenants.PUT("/:id", tenantHandler.UpdateTenant)
		tenants.DELETE("/:id", tenantHandler.DeleteTenant)
	}

	return router
}
