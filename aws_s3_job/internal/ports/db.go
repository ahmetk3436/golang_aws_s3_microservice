package ports

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/application/core/domain"
	"context"
)

type DBPort interface {
	SaveProductsWithTransaction(ctx context.Context, products []*domain.Product) error
}
