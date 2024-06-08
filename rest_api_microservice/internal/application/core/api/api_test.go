package api

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/application/core/domain"
	"context"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDBPort struct {
	mock.Mock
}

func (m *MockDBPort) GetProductById(id int) (*domain.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func TestGetProductById(t *testing.T) {
	mockDB := new(MockDBPort)

	expectedProduct := &domain.Product{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title: "Test Product",
	}
	ctx := context.Background()

	mockDB.On("GetProductById", 1).Return(expectedProduct, nil)
	application := NewApplication(ctx, mockDB)

	product, err := application.GetProductById(1)

	if product == nil || product.ID != expectedProduct.ID || product.Title != expectedProduct.Title {
		t.Errorf("GetProductById returned unexpected result. Expected: %v, Got: %v", expectedProduct, product)
	}

	if err != nil {
		t.Errorf("GetProductById returned an error: %v", err)
	}

	mockDB.AssertCalled(t, "GetProductById", 1)
}
