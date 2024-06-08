package rest_api

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/application/core/domain"
	"fmt"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// GetProductById Get Product By id
// @Summary Get Product By id
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id query int true "Product ID"
// @Success 200 {object} domain.Response "Product retrieved successfully"
// @Failure 400 {object} domain.Response "Bad request"
// @Failure 500 {object} domain.Response "Internal server error"
// @Router /product [get]
func (a *Adapter) GetProductById(c *fiber.Ctx) error {
	log.Info("Getting product...")
	id := c.Query("id")

	product, err := a.api.GetProductById(stringConvertToInt(id))

	if err != nil {
		fmt.Printf("error getting product: %v\n", err)
		return c.Status(500).JSON(domain.Response{
			Status:  500,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.Response{
		Status:  200,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

func stringConvertToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("error converting to int: %v\n", err)
	}
	return i
}
