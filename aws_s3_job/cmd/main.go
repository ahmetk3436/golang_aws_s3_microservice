package main

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/config"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/adapters/aws_s3"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/adapters/db"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/adapters/file"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/aws_s3_job/internal/application/core/api"
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Job's starting...")
	startTime := time.Now()
	appContext := context.Background()
	dbAdapter, err := db.NewAdapter(config.GetDataSourceUrl())
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		panic(err)
	}
	fileAdapter := file.NewAdapter()
	awsAdapter := aws_s3.NewAdapter()
	application := api.NewApplication(dbAdapter, fileAdapter, awsAdapter)

	err = application.ProcessS3Files(appContext, config.GetBucketName())
	if err != nil {
		fmt.Printf("Error processing files: %v\n", err)
		panic(err)
	}
	endTime := time.Now()
	fmt.Printf("Job took %s\n", endTime.Sub(startTime))
}
