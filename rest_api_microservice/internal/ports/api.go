package ports

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/application/core/domain"
)

type Api interface {
	GetProductById(id int) (*domain.Product, error)
}
