package server

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"
	"vladislavsperkanuks/feature-toggles/pkg/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
	*gin.Engine
}

func New(db *gorm.DB) *Server {
	s := &Server{
		db:     db,
		Engine: gin.Default(),
	}

	setCORS(s)
	setHandlers(s)

	return s
}

func setCORS(s *Server) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	s.Use(cors.New(config))
}

func setHandlers(s *Server) {
	// Customers
	s.POST("/api/v1/customers", s.createCustomer)
	s.GET("/api/v1/customers", s.getCustomers)
	// Feature toggles
	s.GET("/api/v1/features", s.getFeatureToggles)
	s.POST("/api/v1/features", s.createFeatureToggle)
	s.PUT("/api/v1/features/:id", s.updateFeatureToggle)
	// Feature requests by customer
	s.POST("/api/v1/customers/:id", s.getCustomerToggles)
}

// Handlers

// POST /api/v1/costumers.
func (s *Server) createCustomer(c *gin.Context) {
	var customer model.Customer

	s.db.Save(&customer)
	c.JSON(http.StatusOK, customer)
}

// GET /api/v1/customers.
func (s *Server) getCustomers(c *gin.Context) {
	var customers []model.Customer

	s.db.Find(&customers)
	c.JSON(http.StatusOK, customers)
}

// GET /api/v1/features.
func (s *Server) getFeatureToggles(c *gin.Context) {
	var toggles []model.FeatureToggle

	s.db.Preload("Customers").Find(&toggles)
	c.JSON(http.StatusOK, toggles)
}

// POST /api/v1/features.
func (s *Server) createFeatureToggle(c *gin.Context) {
	var toggle model.FeatureToggle

	if err := c.BindJSON(&toggle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// Check if feature toggle with the same technical name already exists
	var existingToggle model.FeatureToggle

	s.db.Where("technical_name = ?", toggle.TechnicalName).First(&existingToggle)

	if existingToggle.TechnicalName != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feature toggle with the same technical name already exists"})
		return
	}

	s.db.Save(&toggle)

	c.JSON(http.StatusOK, toggle)
}

// PATCH /api/v1/features/:id.
func (s *Server) updateFeatureToggle(c *gin.Context) {
	var dst model.FeatureToggle
	// First, we find the existing feature toggle
	s.db.Where("id = ?", c.Param("id")).First(&dst)

	if dst.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feature toggle not found"})
		return
	}

	var src model.FeatureToggle

	// HACK: ignore the required fields, because update request might not contain all fields
	if err := c.ShouldBindJSON(&src); err != nil {
		var fieldValidationError validator.FieldError
		if ok := errors.As(err, &fieldValidationError); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	for _, customer := range src.Customers {
		var existingCustomer model.Customer

		s.db.Where("id = ?", customer.ID).First(&existingCustomer)

		if existingCustomer.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add non-existing customer to the feature toggle"})
			return
		}
	}

	model.UpdateFeatureToggle(&src, &dst)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Save the FeatureToggle first (without Customers)
		if err := tx.Save(&dst).Error; err != nil {
			return fmt.Errorf("save feature toggle: %w", err)
		}

		// Replace the Customers
		if err := tx.Model(&dst).Association("Customers").Replace(dst.Customers); err != nil {
			return fmt.Errorf("replace customers: %w", err)
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dst)
}

// POST /api/v1/customers/:id.
func (s *Server) getCustomerToggles(c *gin.Context) {
	var customer model.Customer

	s.db.Where("id = ?", c.Param("id")).First(&customer)

	if customer.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var req struct {
		Features []string `binding:"required" json:"features"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch all toggles
	var toggles []model.FeatureToggle

	s.db.Where("technical_name IN ?", req.Features).Preload("Customers").Find(&toggles)

	type response struct {
		TechnicalName string `json:"technicalName"`
		IsActive      bool   `json:"isActive"`
		IsInverted    bool   `json:"isInverted"`
		IsExpired     bool   `json:"isExpired"`
	}

	resp := make([]response, 0, len(req.Features))

	for _, featureName := range req.Features {
		// Check if feature toggle with the technical name exists in db
		idx := slices.IndexFunc(toggles, func(t model.FeatureToggle) bool { return t.TechnicalName == featureName })

		if idx == -1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Feature toggle with technical name '%s' not found", featureName),
			})

			return
		}

		resp = append(resp, response{
			TechnicalName: toggles[idx].TechnicalName,
			IsActive:      toggles[idx].IsActiveForCustomer(customer),
			IsInverted:    *toggles[idx].IsInverted,
			IsExpired:     toggles[idx].ExpiresOn != nil && time.Now().After(*toggles[idx].ExpiresOn),
		})
	}

	c.JSON(http.StatusOK, resp)
}
