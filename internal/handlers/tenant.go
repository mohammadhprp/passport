package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohammadhprp/passport/internal/services"
)

// TenantHandler wires HTTP endpoints to the tenant service.
type TenantHandler struct {
	service *services.TenantService
}

// NewTenantHandler constructs a tenant handler.
func NewTenantHandler(service *services.TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

type tenantPayload struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// CreateTenant handles POST /tenants.
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var payload tenantPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.service.CreateTenant(c.Request.Context(), payload.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

// ListTenants handles GET /tenants.
func (h *TenantHandler) ListTenants(c *gin.Context) {
	tenants, err := h.service.ListTenants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

// GetTenant handles GET /tenants/:id.
func (h *TenantHandler) GetTenant(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	tenant, err := h.service.GetTenant(c.Request.Context(), id)
	if err != nil {
		switch err {
		case services.ErrTenantNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// UpdateTenant handles PUT /tenants/:id.
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var payload tenantPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.service.UpdateTenant(c.Request.Context(), id, payload.Name)
	if err != nil {
		switch err {
		case services.ErrTenantNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// DeleteTenant handles DELETE /tenants/:id.
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteTenant(c.Request.Context(), id); err != nil {
		switch err {
		case services.ErrTenantNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	raw := c.Param(name)
	id, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": name + " must be a valid UUID"})
		return uuid.Nil, false
	}
	return id, true
}
