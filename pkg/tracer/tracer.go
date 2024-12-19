package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func Init(ctx context.Context, serviceName, jaegerURL string) (trace.Tracer, func(ctx context.Context) error, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpointURL(jaegerURL),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("service.name", serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	tracer := otel.Tracer(serviceName)

	return tracer, tp.Shutdown, nil
}
