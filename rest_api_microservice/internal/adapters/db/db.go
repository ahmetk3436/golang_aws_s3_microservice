package db

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/application/core/domain"
	"context"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Adapter struct {
	ctx context.Context
	db  *gorm.DB
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	var db *gorm.DB
	var err error
	retries := 3

	for i := 1; i <= retries; i++ {
		db, err = gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
		if err == nil {
			return &Adapter{db: db}, nil
		}
		fmt.Printf("Database connection error (attempt %d): %v\n", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	if db != nil {
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}
	}

	return nil, fmt.Errorf("could not establish database connection after %d retries: %v", retries, err)
}

func (a *Adapter) GetProductById(id int) (*domain.Product, error) {
	var product *domain.Product
	var err error
	retries := 3

	for i := 0; i < retries; i++ {
		if err = a.db.WithContext(a.ctx).First(&product, id).Error; err == nil {
			return product, nil
		}
		fmt.Printf("Error finding product %v (attempt %d)\n", err, i+1)
		time.Sleep(time.Second * time.Duration(i+1))
	}

	return nil, fmt.Errorf("could not find product after %d retries: %v", retries, err)
}
