package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"

	"github.com/redrru/fantasy-dota/pkg/log"
	"github.com/redrru/fantasy-dota/pkg/tracing"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
	}
}

func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	ctx, span := tracing.DefaultTracer().Start(ctx, "HttpClient")
	defer span.End()

	ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	log.GetLogger().Debug(ctx, "Sending GET request", zap.String("url", url))
	res, err := c.client.Do(req)
	defer func() {
		if res == nil || res.Body == nil {
			return
		}
		if err := res.Body.Close(); err != nil {
			log.GetLogger().Warn(ctx, "Close http body", zap.String("url", url), zap.Error(err))
		}
	}()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: '%v', body: '%s'", res.Status, string(body))
	}

	return body, nil
}
