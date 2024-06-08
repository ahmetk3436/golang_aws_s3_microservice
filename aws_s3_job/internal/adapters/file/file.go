package file

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/application/core/domain"
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type Adapter struct {
}

func NewAdapter() *Adapter {
	return &Adapter{}
}
func (a *Adapter) DeleteFile(filePath string) {
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("Error deleting file '%s': %v\n", filePath, err)
	} else {
		fmt.Printf("File '%s' deleted successfully\n", filePath)
	}
}
func (a *Adapter) ReadFile(ctx context.Context, filePath string) ([]*domain.Product, error) {
	done := ctx.Done()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}(file)

	decoder := json.NewDecoder(file)

	var products []*domain.Product
	for decoder.More() {
		select {
		case <-done:
			fmt.Println("File processing cancelled due to context cancellation")
			return nil, ctx.Err()
		default:
			var product domain.Product
			if err := decoder.Decode(&product); err != nil {
				fmt.Printf("Error decoding file: %s\n", err)
				return nil, err
			}
			products = append(products, &product)
		}
	}

	return products, nil
}
