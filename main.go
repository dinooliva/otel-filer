package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
    "go.opentelemetry.io/otel/sdk/resource"
    sdkmetric "go.opentelemetry.io/otel/sdk/metric"
    semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func main() {
    ctx := context.Background()

    // Set up the OpenTelemetry Metric exporter
    exporter, err := otlpmetrichttp.New(ctx)
    if err != nil {
        log.Fatalf("failed to create OTLP metric exporter: %v", err)
    }

    // Configure the meter provider
    meterProvider := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(
            sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(2*time.Second)),
        ),
        sdkmetric.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("otel-example-app"),
        )),
    )

    otel.SetMeterProvider(meterProvider)
    defer func() {
        if err := meterProvider.Shutdown(ctx); err != nil {
            log.Fatalf("failed to shutdown MeterProvider: %v", err)
        }
    }()

    meter := otel.Meter("otel-example-meter")
    requestCount, err := meter.Int64Counter("request_count")
    if err != nil {
        log.Fatalf("failed to create request_count counter: %v", err)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Increment the request count
        requestCount.Add(ctx, 1)

        fmt.Fprintf(w, "Hello, OpenTelemetry Metrics!")
        log.Println("Request handled")
    })

    log.Println("Server is running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
