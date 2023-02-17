package main

import (
	"context"
	_ "embed"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

//go:embed otel-collector-config.yaml
var configYAML []byte

func run(ctx context.Context) error {
	errorCh := make(chan error)

	factories, err := components()
	if err != nil {
		return err
	}

	uriLocation := "yaml:" + string(configYAML)
	provider := yamlprovider.New()
	set := otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs:      []string{uriLocation},
			Providers: map[string]confmap.Provider{provider.Scheme(): provider},
		},
	}

	cp, err := otelcol.NewConfigProvider(set)
	if err != nil {
		return err
	}

	cfg, err := cp.Get(context.Background(), factories)
	if err != nil {
		return err
	}

	svc, err := service.New(ctx, service.Settings{
		BuildInfo: component.BuildInfo{
			Command:     "builtin",
			Description: "Builtin OpenTelemetry Collector",
			Version:     "0.0.0",
		},
		Receivers:         receiver.NewBuilder(cfg.Receivers, factories.Receivers),
		Processors:        processor.NewBuilder(cfg.Processors, factories.Processors),
		Exporters:         exporter.NewBuilder(cfg.Exporters, factories.Exporters),
		Connectors:        connector.NewBuilder(cfg.Connectors, factories.Connectors),
		Extensions:        extension.NewBuilder(cfg.Extensions, factories.Extensions),
		AsyncErrorChannel: errorCh,
		LoggingOptions:    []zap.Option{},
	}, cfg.Service)
	if err != nil {
		return err
	}

	err = svc.Start(ctx)
	if err != nil {
		return err
	}

	select {
	case e := <-errorCh:
		return e
	case <-ctx.Done():
	}
	return nil
}

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
