package api

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/application/core/domain"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/ports"
	"context"
	"fmt"
)

type Application struct {
	db  ports.DBPort
	ctx context.Context
}

func NewApplication(ctx context.Context, db ports.DBPort) *Application {
	return &Application{
		db:  db,
		ctx: ctx,
	}
}

func (a *Application) GetProductById(id int) (*domain.Product, error) {
	product, err := a.db.GetProductById(id)
	if err != nil {
		fmt.Printf("Error getting product by id: %v", err)
		return nil, err
	}
	return product, nil
}
