package api

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/application/core/domain"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/ports"
	"context"
	"fmt"
	"path/filepath"
)

type Application struct {
	db      ports.DBPort
	fileJob ports.FileJobPort
	aws     ports.Aws
}

func NewApplication(db ports.DBPort, fileJob ports.FileJobPort, aws ports.Aws) *Application {
	return &Application{
		db:      db,
		fileJob: fileJob,
		aws:     aws,
	}
}

func (a *Application) ProcessS3Files(ctx context.Context, bucketName string) error {
	sess, err := a.aws.LoginAws()
	if err != nil {
		return err
	}

	fileDataChan := make(chan *domain.FileData)

	go func() {
		defer close(fileDataChan)
		err := a.aws.Download(sess, bucketName, fileDataChan)
		if err != nil {
			fmt.Printf("Error downloading files from S3: %v\n", err)
		}
	}()

	var productsToSave []*domain.Product

	for fileData := range fileDataChan {
		if fileData.Success {
			fmt.Printf("File '%s' successfully downloaded\n", *fileData.Key)
			filePath := filepath.Join("./downloads", *fileData.Key)

			products, err := a.fileJob.ReadFile(ctx, filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", filePath, err)
				continue
			}

			productsToSave = append(productsToSave, products...)
		} else {
			fmt.Printf("Error downloading file '%s': %v\n", *fileData.Key, fileData.Error)
		}
	}

	if err := a.db.SaveProductsWithTransaction(ctx, productsToSave); err != nil {
		fmt.Printf("Error saving products with transaction: %v\n", err)
		return err
	}

	return nil
}
