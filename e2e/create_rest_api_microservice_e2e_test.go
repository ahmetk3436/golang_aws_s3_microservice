package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	tcCompose "github.com/testcontainers/testcontainers-go/modules/compose"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

type CreateOrderTestSuite struct {
	suite.Suite
	compose    *tcCompose.LocalDockerCompose
	identifier string
}

func (c *CreateOrderTestSuite) SetupSuite() {
	composeFilePaths := []string{"resources/docker-compose.yml"}

	c.identifier = strings.ToLower(uuid.New().String())

	dockerCompose := tcCompose.NewLocalDockerCompose(composeFilePaths, c.identifier)
	c.compose = dockerCompose

	execError := dockerCompose.
		WithCommand([]string{"up", "-d"}).
		Invoke()
	if execError.Error != nil {
		log.Fatalf("Could not run compose stack: %v", execError)
	}
}

func (c *CreateOrderTestSuite) Test_Should_Get_Response_From_DB() {

	request, err := http.Get("http://localhost:9999/product?id=3")
	if err != nil {
		fmt.Printf("Error getting products: %v\n", err)
		c.T().Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing body: %v\n", err)
		}
	}(request.Body)

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		c.T().Fatal(err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error unmarshalling body: %v\n", err)
		c.T().Fatal(err)
	}

	c.Equal(float32(15811.11), response.Data.Price)
	c.Equal("title3", response.Data.Title)
}

func (c *CreateOrderTestSuite) TearDownSuite() {
	execError := c.compose.Down()
	if execError.Error != nil {
		log.Fatalf("Could not shutdown compose stack: %v", execError)
	}
}

func TestCreateOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}

type Response struct {
	Message string  `json:"message"`
	Status  int     `json:"status"`
	Data    Product `json:"data"`
}

type Product struct {
	ID    int     `json:"ID"`
	Title string  `json:"title"`
	Price float32 `json:"price"`
}
