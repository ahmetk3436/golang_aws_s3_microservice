package db

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/application/core/domain"
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

type OrderDatabaseTestSuite struct {
	suite.Suite
	DataSourceUrl string
}

func (o *OrderDatabaseTestSuite) SetupSuite() {
	ctx := context.Background()
	port := "3306/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("root:golang-aws-s3-backend-task@tcp(%s:%s)/products?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())
	}
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:8.0.30",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "golang-aws-s3-backend-task",
			"MYSQL_DATABASE":      "products",
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "mysql", dbURL).WithStartupTimeout(time.Second * 30),
	}
	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal("Failed to start MySQL container:", err)
	}
	endpoint, _ := mysqlContainer.Endpoint(ctx, "")
	o.DataSourceUrl = fmt.Sprintf("root:golang-aws-s3-backend-task@tcp(%s)/products?charset=utf8mb4&parseTime=True&loc=Local", endpoint)
}

func (o *OrderDatabaseTestSuite) Test_Should_Get_Order() {
	adapter, _ := NewAdapter(o.DataSourceUrl)
	ctx := context.Background()
	productId := 1
	product, err := adapter.GetProductById(ctx, productId)
	if err != nil {
		fmt.Printf("Failed to get product by id: %v\n", productId)
	}
	var nilProduct *domain.Product
	o.Equal(nilProduct, product)
}

func TestOrderDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDatabaseTestSuite))
}
