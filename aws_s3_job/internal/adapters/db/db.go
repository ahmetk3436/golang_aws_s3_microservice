package db

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/application/core/domain"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Adapter struct {
	db *gorm.DB
}

func (a *Adapter) SaveProductsWithTransaction(ctx context.Context, products []*domain.Product) error {
	tx := a.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			fmt.Printf("failed to commit transaction: %v\n", err)
		}
	}()

	productMap := make(map[int64]*domain.Product)
	for _, product := range products {
		productMap[int64(product.ID)] = product
	}

	for _, product := range productMap {
		a.saveProductInTx(ctx, tx, product)
	}

	return nil
}

func (a *Adapter) saveProductInTx(ctx context.Context, tx *gorm.DB, product *domain.Product) {
	err := tx.WithContext(ctx).Create(product).Error
	if err != nil {
		fmt.Printf("error while saving product: %v", err)
	}
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	var db *gorm.DB
	var err error
	const maxRetries = 3
	const retryInterval = 10 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
		if err == nil {
			fmt.Printf("Database connection successful after %d attempts\n", attempt)
			break
		}

		fmt.Printf("Error connecting to database (attempt %d/%d): %v\n", attempt, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}
	if db != nil {
		err = db.AutoMigrate(&domain.Product{})
		if err != nil {
			return nil, fmt.Errorf("failed to migrate database: %v", err)
		}
	}

	return &Adapter{db: db}, nil
}
