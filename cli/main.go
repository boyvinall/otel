package main

import (
	"context"
	"log"

	// "go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

func run() error {
	ctx := context.Background()

	grpcExpOpt := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint("127.0.0.1:4317"),
		otlpmetricgrpc.WithDialOption(
			grpc.WithBlock(),
		),
		otlpmetricgrpc.WithInsecure(),
	}
	exp, err := otlpmetricgrpc.New(ctx, grpcExpOpt...)
	if err != nil {
		return err
	}
	defer func() {
		if e := exp.ForceFlush(ctx); e != nil {
			log.Println("failed to flush exporter", e)
		}
		if e := exp.Shutdown(ctx); e != nil {
			log.Println("failed to stop the exporter", e)
		}
	}()

	// This reader is used as a stand-in for a reader that will actually export
	// data. See exporters in the go.opentelemetry.io/otel/exporters package
	// for more information.
	reader := metric.NewManualReader()

	// See the go.opentelemetry.io/otel/sdk/resource package for more
	// information about how to create and use Resources.
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("my-service"),
		semconv.ServiceVersionKey.String("v0.1.0"),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)
	// global.SetMeterProvider(meterProvider)
	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Println("failed to stop the meter provider", err)
		}
	}()
	// The MeterProvider is configured and registered globally. You can now run
	// your code instrumented with the OpenTelemetry API that uses the global
	// MeterProvider without having to pass this MeterProvider instance. Or,
	// you can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.
	//
	// See the go.opentelemetry.io/otel/metric package for more information
	// about the metric API.

	m := meterProvider.Meter("app_or_package_name")
	foo, err := m.Int64Counter("foo",
		instrument.WithDescription("fooo"),
		instrument.WithUnit(unit.Milliseconds))
	foo.Add(ctx, 1)
	foo.Add(ctx, 1)
	foo.Add(ctx, 1)

	updown, err := m.Int64UpDownCounter("updown", instrument.WithDescription("uppy downy description"))
	if err != nil {
		return err
	}
	updown.Add(ctx, 1)
	updown.Add(ctx, 1)
	updown.Add(ctx, 1)

	rm, err := reader.Collect(ctx)
	if err != nil {
		return err
	}

	err = exp.Export(ctx, rm)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
