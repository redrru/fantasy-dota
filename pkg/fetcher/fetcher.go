package apifetcher

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/redrru/fantasy-dota/pkg/http"
	"github.com/redrru/fantasy-dota/pkg/log"
	"github.com/redrru/fantasy-dota/pkg/tracing"
)

type Handler interface {
	Handle(ctx context.Context, response []byte) error
	GetRefreshTime() time.Duration
	GetURL() string
}

type Fetcher struct {
	httpClient *http.Client
	handlers   []Handler

	handler      chan Handler
	tickersClose []chan struct{}
	close        chan struct{}
}

func NewFetcher() *Fetcher {
	fetcher := &Fetcher{
		httpClient: http.NewClient(),
		handler:    make(chan Handler, 5),
		close:      make(chan struct{}),
	}

	return fetcher
}

func (f *Fetcher) RegisterHandlers(handlers ...Handler) {
	f.handlers = append(f.handlers, handlers...)
}

func (f *Fetcher) Run() {
	f.initTickers()

	for {
		select {
		case <-f.close:
			return
		case handler := <-f.handler:
			ctx := context.Background()
			logger := log.GetLogger().With(zap.String("url", handler.GetURL()))

			f.withTracing(ctx, func(ctx context.Context) (err error) {
				defer func() {
					if r := recover(); r != nil {
						if panicErr, ok := r.(error); ok {
							err = panicErr
							logger.Error(ctx, "Fetch panic", zap.Error(panicErr))
						}
					}
				}()

				resp, err := f.fetch(ctx, handler)
				if err != nil {
					logger.Error(ctx, "Fetch", zap.Error(err))
					return err
				}

				err = handler.Handle(ctx, resp)
				if err != nil {
					logger.Error(ctx, "Handle", zap.Error(err))
					return err
				}

				return err
			})
		}
	}
}

func (f *Fetcher) Close() error {
	f.close <- struct{}{}

	for _, ch := range f.tickersClose {
		ch <- struct{}{}
	}

	log.GetLogger().Debug(context.Background(), "Fetcher exited")

	return nil
}

func (f *Fetcher) initTickers() {
	for _, handler := range f.handlers {
		handler := handler
		stop := make(chan struct{})
		f.tickersClose = append(f.tickersClose, stop)

		go func() {
			ticker := time.NewTicker(handler.GetRefreshTime())
			defer ticker.Stop()

			for {
				select {
				case <-stop:
					return
				case <-ticker.C:
					f.handler <- handler
				}
			}
		}()
	}
}

func (f *Fetcher) fetch(ctx context.Context, handler Handler) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	return f.httpClient.Get(ctx, handler.GetURL())
}

func (f *Fetcher) withTracing(ctx context.Context, do func(ctx context.Context) error) {
	ctx, span := tracing.DefaultTracer().Start(ctx, "FetchAPI")
	defer span.End()

	if err := do(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
