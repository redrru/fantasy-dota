package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.uber.org/zap"

	"github.com/redrru/fantasy-dota/pkg/env"
	apifetcher "github.com/redrru/fantasy-dota/pkg/fetcher"
	"github.com/redrru/fantasy-dota/pkg/log"
	"github.com/redrru/fantasy-dota/pkg/middleware"
)

type Closer func() error

type Application struct {
	name string
	tp   *trace.TracerProvider

	env      env.Env
	shutdown chan os.Signal
	err      chan error

	fetcher *apifetcher.Fetcher
	http    *echo.Echo
	closers []Closer
}

func NewApplication() *Application {
	app := &Application{
		name:     "fantasy-dota",
		env:      env.GetEnv(),
		shutdown: make(chan os.Signal, 1),
		err:      make(chan error, 1),
		fetcher:  apifetcher.NewFetcher(),
	}

	app.closers = append(app.closers, app.fetcher.Close)

	app.initTracing()

	return app
}

func (a *Application) RegisterFetchers(handlers ...apifetcher.Handler) {
	a.fetcher.RegisterHandlers(handlers...)
}

func (a *Application) RegisterHTTP(e *echo.Echo) {
	a.http = e
}

func (a *Application) Run() {
	defer log.GetLogger().Info(context.Background(), "App exited")

	signal.Notify(a.shutdown, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	go a.fetcher.Run()
	go a.serverHTTP()

	select {
	case sig := <-a.shutdown:
		log.GetLogger().Info(context.Background(), "Got signal", zap.String("signal", sig.String()))
	case err := <-a.err:
		log.GetLogger().Info(context.Background(), "Got fatal err", zap.Error(err))
	}

	a.close()
}

func (a *Application) serverHTTP() {
	a.http.Use(
		middleware.TracingMiddleware(a.name),
		middleware.LoggingMiddleware(),
		middleware.RecoveringMiddleware(),
	)

	if err := a.http.Start(":8080"); err != nil {
		a.err <- err
	}
}

func (a *Application) close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.tp.Shutdown(ctx); err != nil {
		log.GetLogger().Error(context.Background(), "TracerProvider shutdown error", zap.Error(err))
	}

	for _, closer := range a.closers {
		if err := closer(); err != nil {
			log.GetLogger().Error(context.Background(), "Shutdown error", zap.Error(err))
		}
	}

	_ = log.GetLogger().Sync()
}

const (
	jaegerHost = "JAEGER_AGENT_HOST"
	jaegerPort = "JAEGER_AGENT_PORT"
)

func (a *Application) initTracing() {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(a.env[jaegerHost]),
		jaeger.WithAgentPort(a.env[jaegerPort]),
	))
	if err != nil {
		panic(err)
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(a.name),
			semconv.ServiceVersionKey.String("v0.0.1"),
		),
	)
	if err != nil {
		panic(err)
	}

	a.tp = trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(r),
	)

	otel.SetTracerProvider(a.tp)
}
