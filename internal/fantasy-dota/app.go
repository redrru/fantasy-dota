package application

import (
	"context"
	"fmt"
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

	postgres "github.com/redrru/fantasy-dota/pkg/db"
	"github.com/redrru/fantasy-dota/pkg/env"
	httpfetcher "github.com/redrru/fantasy-dota/pkg/fetcher"
	"github.com/redrru/fantasy-dota/pkg/log"
	"github.com/redrru/fantasy-dota/pkg/middleware"
)

const (
	appVersionEnv = "APP_VERSION"

	httpPortEnv = "APP_HTTP_PORT"

	jaegerHostEnv = "JAEGER_AGENT_HOST"
	jaegerPortEnv = "JAEGER_AGENT_PORT"

	postgresDSN             = "PG_DSN"
	postgresMaxIdleConns    = "PG_MAX_IDLE_CONNS"
	postgresMaxOpenConns    = "PG_MAX_OPEN_CONNS"
	postgresConnMaxLifetime = "PG_CONN_MAX_LIFETIME"
	postgresConnMaxIdleTime = "PG_CONN_MAX_IDLE_TIME"

	logStr = "[APP] %s"
)

type Closer func() error

type Application struct {
	name string
	env  env.Env

	fetcher *httpfetcher.Fetcher
	http    *echo.Echo
	DB      *postgres.DB
	tp      *trace.TracerProvider

	closers  []Closer
	dbModels []interface{}

	shutdown chan os.Signal
	httpErr  chan error
}

func NewApplication() *Application {
	app := &Application{
		name:     "fantasy-dota",
		env:      env.GetEnv(),
		shutdown: make(chan os.Signal, 1),
		httpErr:  make(chan error, 1),
		fetcher:  httpfetcher.NewFetcher(),
	}

	app.closers = append(app.closers, app.fetcher.Close)

	app.initTracing()
	app.initDB()

	return app
}

func (a *Application) RegisterFetchers(handlers ...httpfetcher.Handler) {
	a.fetcher.RegisterHandlers(handlers...)
}

func (a *Application) RegisterHTTP(e *echo.Echo) {
	a.http = e
}

func (a *Application) RegisterMigrationModel(models ...interface{}) {
	a.dbModels = append(a.dbModels, models...)
}

func (a *Application) Run() {
	defer log.GetLogger().Info(context.Background(), fmt.Sprintf(logStr, "Exited"))

	signal.Notify(a.shutdown, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	a.migrationDB()

	go a.fetcher.Run()
	go a.serverHTTP()

	log.GetLogger().Info(context.Background(), fmt.Sprintf(logStr, "Started"))

	select {
	case sig := <-a.shutdown:
		log.GetLogger().Info(context.Background(), fmt.Sprintf(logStr, "Got signal"), zap.String("signal", sig.String()))
	case err := <-a.httpErr:
		log.GetLogger().Info(context.Background(), fmt.Sprintf(logStr, "Got http fatal err"), zap.Error(err))
	}

	a.stop()
}

func (a *Application) initDB() {
	cfg := postgres.Config{
		DSN:             a.env.GetString(postgresDSN),
		MaxIdleConns:    a.env.GetInt(postgresMaxIdleConns),
		MaxOpenConns:    a.env.GetInt(postgresMaxOpenConns),
		ConnMaxLifetime: a.env.GetDuration(postgresConnMaxLifetime),
		ConnMaxIdleTime: a.env.GetDuration(postgresConnMaxIdleTime),
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	a.DB = db
	a.closers = append(a.closers, a.DB.Close)

	a.waitDB()
}

func (a *Application) waitDB() {
	log.GetLogger().Info(context.Background(), fmt.Sprintf(logStr, "Waiting DB up..."))

	check := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := a.DB.Ping(ctx)
		if err != nil {
			log.GetLogger().Error(context.Background(), fmt.Sprintf(logStr, "Waiting DB up"), zap.Error(err))
		} else {
			log.GetLogger().Info(context.Background(), fmt.Sprintf(logStr, "DB up"))
		}
		return err
	}

	err := check()
	if err == nil {
		return
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for i := 0; i < 5; i++ {
		<-ticker.C
		if err = check(); err == nil {
			return
		}
	}
	if err != nil {
		panic(fmt.Errorf("wait DB up timeout"))
	}
}

func (a *Application) migrationDB() {
	if err := a.DB.Gorm.AutoMigrate(a.dbModels...); err != nil {
		panic(fmt.Errorf("db migration failed: %w", err))
	}
}

func (a *Application) serverHTTP() {
	a.http.Use(
		middleware.TracingMiddleware(a.name),
		middleware.LoggingMiddleware(),
		middleware.RecoveringMiddleware(),
	)

	if err := a.http.Start(fmt.Sprintf(":%s", a.env.GetString(httpPortEnv))); err != nil {
		a.httpErr <- err
	}
}

func (a *Application) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.tp.Shutdown(ctx); err != nil {
		log.GetLogger().Error(context.Background(), fmt.Sprintf(logStr, "TracerProvider shutdown error"), zap.Error(err))
	}

	for _, closer := range a.closers {
		if err := closer(); err != nil {
			log.GetLogger().Error(context.Background(), fmt.Sprintf(logStr, "Shutdown error"), zap.Error(err))
		}
	}

	_ = log.GetLogger().Sync()
}

func (a *Application) initTracing() {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(a.env.GetString(jaegerHostEnv)),
		jaeger.WithAgentPort(a.env.GetString(jaegerPortEnv)),
	))
	if err != nil {
		panic(err)
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(a.name),
			semconv.ServiceVersionKey.String(a.env.GetString(appVersionEnv)),
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
