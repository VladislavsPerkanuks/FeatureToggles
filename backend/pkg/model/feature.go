package model

import (
	"slices"
	"time"

	"gorm.io/gorm"
)

type FeatureToggle struct {
	gorm.Model
	DisplayName   *string `json:"displayName"`
	TechnicalName string  `binding:"required" json:"technicalName"`
	// HACK: Pointer as zero value fails the required validation
	IsInverted  *bool      `binding:"required"                        json:"isInverted"`
	ExpiresOn   *time.Time `json:"expiresOn"`
	Description *string    `json:"description"`
	IsArchived  bool       `json:"isArchived"`
	Customers   []Customer `gorm:"many2many:featuretoggle_customers;" json:"customers"`
}

// IsActiveForCustomer checks if the feature toggle is active for the given customer.
func (f *FeatureToggle) IsActiveForCustomer(customer Customer) bool {
	var isActive bool

	// customer is in the list of the feature toggle:
	if slices.ContainsFunc(f.Customers, func(c Customer) bool { return c.ID == customer.ID }) {
		isActive = !isActive
	}

	// if archived or expired, the feature is not active:
	if f.IsArchived {
		isActive = false
	}

	if f.ExpiresOn != nil && time.Now().After(*f.ExpiresOn) {
		isActive = !isActive
	}

	// if inverted status is flipped
	if f.IsInverted != nil && *f.IsInverted {
		isActive = !isActive
	}

	return isActive
}

// UpdateFeatureToggle updates the dst FeatureToggle with the src FeatureToggle fields.
func UpdateFeatureToggle(src, dst *FeatureToggle) {
	if src == nil || dst == nil {
		return
	}

	dst.DisplayName = src.DisplayName
	dst.TechnicalName = src.TechnicalName
	dst.IsInverted = src.IsInverted
	dst.ExpiresOn = src.ExpiresOn
	dst.Description = src.Description
	dst.IsArchived = src.IsArchived
	dst.Customers = src.Customers
}
