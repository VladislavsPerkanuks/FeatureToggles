package db

import (
	"fmt"
	"os"
	"time"
	"vladislavsperkanuks/feature-toggles/pkg/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// initData initializes the database with some data for demonstration purposes.
func initData(db *gorm.DB) {
	// if already initialized, skip
	if db.First(&model.Customer{}).Error == nil {
		return
	}

	// Create a new customers
	customer1 := model.Customer{}
	customer2 := model.Customer{}
	customer3 := model.Customer{}
	customer4 := model.Customer{}

	db.Create(&customer1)
	db.Create(&customer2)
	db.Create(&customer3)
	db.Create(&customer4)

	// Create a features
	allFields := model.FeatureToggle{
		DisplayName:   ptr("All fields feature"),
		TechnicalName: "Feature with all possible fields",
		IsInverted:    ptr(true),
		ExpiresOn:     ptr(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		Description:   ptr("This feature has all possible fields"),
		IsArchived:    false,
		Customers:     []model.Customer{customer1, customer2},
	}

	onlyNecessaryFields := model.FeatureToggle{
		TechnicalName: "Feature with only necessary fields",
		IsInverted:    ptr(false),
		Customers:     []model.Customer{customer3},
	}

	allCustomersFeature := model.FeatureToggle{
		TechnicalName: "Feature for all customers",
		IsInverted:    ptr(false),
		Customers:     []model.Customer{customer1, customer2, customer3, customer4},
		DisplayName:   ptr("Feature for all customers"),
	}

	archivedFeature := model.FeatureToggle{
		TechnicalName: "Archived feature",
		IsInverted:    ptr(false),
		IsArchived:    true,
	}

	db.Create(&allFields)
	db.Create(&onlyNecessaryFields)
	db.Create(&allCustomersFeature)
	db.Create(&archivedFeature)
}

func ptr[T any](t T) *T {
	return &t
}

// New creates a new gorm.DB instance and performs auto migrations.
func New(path string) (*gorm.DB, func() error, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.AutoMigrate(&model.FeatureToggle{}, &model.Customer{}); err != nil {
		return nil, nil, fmt.Errorf("auto migrate: %w", err)
	}

	if os.Getenv("DEMO") == "true" {
		initData(db)
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get db instance: %w", err)
	}

	return db, dbInstance.Close, nil
}
