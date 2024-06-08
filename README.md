# Golang Aws S3 Microservice

## What I Used ?
- Hexagonal Architecture
- Unit, Integration and End-To-End Tests
- Dockerfile & Docker-Compose
- Logrus for Aggregate Logs
- Jaeger and Prometheus Implementation With OpenTelemetry
- Swagger 2.0 implementation
- Fiber
- MYSQL
- GORM

## Implementation

- Create Docker Volume And Instances
- Export Env Variables
- Run Aws-S3-Job
- Run Rest-Api-Microservice
- Test Is Endpoint Working

## Create Volume Folder For Mysql (if you want)
```bash
mkdir data
```

## Create Mysql Docker Instance With Volume
```bash
sudo docker run -d --name=golang-aws-s3-mysql -p 3306:3306 -v $pwd/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=golang-aws-s3-microservice mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

## Create Mysql Docker Instance Without Volume
```bash
sudo docker run -d --name=golang-aws-s3-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=golang-aws-s3-microservice mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

## Create golang-aws-s3 Database For Mysql
```bash
sudo docker exec golang-aws-s3-mysql mysql -uroot -pgolang-aws-s3-microservice -e "CREATE DATABASE IF NOT EXISTS golang-aws-s3;"
```

## Create Jaeger All-In-One Docker Instance 

```bash
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest
```

## Export Env Variables
```bash
unset CURRENT_PATH
export CURRENT_PATH=$(pwd)
export DATA_SOURCE_URL="root:golang-aws-s3-microservice@tcp(localhost:3306)/golang-aws-s3?charset=utf8mb4&parseTime=True&loc=Local"
export AWS_REGION="AWS_REGION"
export AWS_ACCESS_KEY_ID="AWS_ACCESS_KEY_ID"
export AWS_BUCKET_NAME="AWS_BUCKET_NAME"
export AWS_SECRET_ACCESS_KEY="AWS_SECRET_ACCESS_KEY"
export DATA_JAEGER_URL="http://localhost:14268/api/traces"
```

## Run Aws S3 Job

```bash
cd $CURRENT_PATH/aws_s3_job/cmd && go run main.go
```

## Run Rest-Api-Microservice

```bash
cd $CURRENT_PATH/rest_api_microservice/cmd && go run main.go
```

## Test Endpoint With Curl

```bash
curl -v "http://localhost:9999/product?id=3"
```

## Run E2E Test

```bash
cd $CURRENT_PATH && cd e2e
go test create_rest_api_microservice_e2e_test.go 
```

# NOTE
***"If any command is not working, ensure that the environment variables are set correctly."***