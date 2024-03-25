package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Test_IsActiveForCustomer tests the IsActiveForCustomer.
func Test_IsActiveForCustomer(t *testing.T) {
	t.Parallel()

	// Prepare the test data
	testCustomer := Customer{Model: gorm.Model{ID: 1}}

	// Features
	tests := []struct {
		name     string
		feature  FeatureToggle
		expected bool
	}{
		{
			// customer is in the list of the feature toggle:
			name: "my-feature-a",
			feature: FeatureToggle{
				TechnicalName: "my-feature-a",
				IsInverted:    ptr(false),
				ExpiresOn:     ptr(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)),
				Customers:     []Customer{{Model: gorm.Model{ID: 1}}},
			},
			expected: true,
		},
		{
			// customer is in the list of the feature toggle, but toggle is inverted:
			name: "my-feature-b",
			feature: FeatureToggle{
				TechnicalName: "my-feature-b",
				IsInverted:    ptr(true),
				ExpiresOn:     ptr(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)),
				Customers:     []Customer{{Model: gorm.Model{ID: 1}}},
			},
			expected: false,
		},
		{
			// customer is NOT in the list of the feature:
			name: "my-feature-c",
			feature: FeatureToggle{
				TechnicalName: "my-feature-c",
				IsInverted:    ptr(false),
				ExpiresOn:     ptr(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expected: false,
		},
		{
			// customer is NOT in the list of the feature, but feature toggle expired:
			name: "my-feature-d",
			feature: FeatureToggle{
				TechnicalName: "my-feature-d",
				IsInverted:    ptr(false),
				ExpiresOn:     ptr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isActive := tt.feature.IsActiveForCustomer(testCustomer)
			require.Equal(t, tt.expected, isActive)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
