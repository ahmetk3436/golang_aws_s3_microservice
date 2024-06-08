package main

import (
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/config"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/adapters/db"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/adapters/rest_api"
	"bitbucket.org/ahmetk3436/golang_aws_s3_microservice/rest_api_microservice/internal/application/core/api"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

const (
	service     = "rest_api_microservice"
	environment = "dev"
	id          = 1
)

type customLogger struct {
	formatter log.JSONFormatter
}

func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	entry.Data["Context"] = span.SpanContext()
	return l.formatter.Format(entry)
}
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

func init() {
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{FieldMap: log.FieldMap{
			"msg": "message",
		}},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
func main() {
	fmt.Println("Rest Api Microservice is starting...")
	ctx := context.Background()
	tp, err := tracerProvider(config.GetJaegerUrl())
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	provider.Meter("api")

	dbAdapter, err := db.NewAdapter(config.GetDataSourceUrl())
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		panic(err)
	}

	application := api.NewApplication(ctx, dbAdapter)

	restApiAdapter := rest_api.NewAdapter(application, 9999)

	err = restApiAdapter.Start()
	if err != nil {
		fmt.Printf("Error starting rest api: %v\n", err)
		panic(err)
	}
}
