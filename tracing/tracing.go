package tracing
import (
    "context"
    "os"
	"google.golang.org/grpc/credentials"
    "log"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func newExporter(ctx context.Context) (*otlptrace.Exporter,error) {
    // Configuration to export data to Honeycomb:
    //
    // 1. The Honeycomb endpoint
    // 2. Your API key, set as the x-honeycomb-team header
    // 3. A Dataset name, set as the x-honeycomb-dataset header
    //
    // A Dataset is a bucket where data gets stored in Honeycomb.
    opts := []otlptracegrpc.Option{
        otlptracegrpc.WithEndpoint("api.honeycomb.io:443"),
        otlptracegrpc.WithHeaders(map[string]string{
            "x-honeycomb-team":    os.Getenv("HONEYCOMB_API_KEY"),
            "x-honeycomb-dataset": os.Getenv("HONEYCOMB_DATASET"),
        }),
        otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
    }

    

    client := otlptracegrpc.NewClient(opts...)
    return otlptrace.New(ctx, client)
}


func newTraceProvider(exp *otlptrace.Exporter)*sdktrace.TracerProvider {
    // The service.name attribute is required.
    resource :=
        resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(os.Getenv("SERVICE_NAME")),
        )

    return sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(resource),
    )
}

func InitTracing(ctx context.Context) *sdktrace.TracerProvider {
	exp, err := newExporter(ctx)

	//exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
        log.Fatalf("failed to initialize exporter: %v", err)
    }
	tp := newTraceProvider(exp)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
        propagation.NewCompositeTextMapPropagator(
            propagation.TraceContext{},
            propagation.Baggage{},
        ),
    )
    return tp
}