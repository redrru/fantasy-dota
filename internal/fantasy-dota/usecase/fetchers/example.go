package fetchers

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/redrru/fantasy-dota/pkg/log"
	"github.com/redrru/fantasy-dota/pkg/tracing"
)

const (
	url               = "http://google.com"
	refreshTimeoutSec = 10
)

type Example struct {
	url  string
	time time.Duration
}

func NewExample() *Example {
	return &Example{
		url:  url,
		time: refreshTimeoutSec * time.Second,
	}
}

func (e *Example) Handle(ctx context.Context, response []byte) error {
	ctx, span := tracing.DefaultTracer().Start(ctx, "ExampleHandle")
	defer span.End()

	log.GetLogger().Debug(ctx, "Handle response", zap.String("url", e.url))

	return nil
}

func (e Example) GetRefreshTime() time.Duration {
	return e.time
}

func (e Example) GetURL() string {
	return e.url
}
