package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/redrru/fantasy-dota/pkg/tracing"
)

func TracingMiddleware(service string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()

			ctx := otel.GetTextMapPropagator().Extract(request.Context(), propagation.HeaderCarrier(request.Header))
			opts := []trace.SpanStartOption{
				trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
				trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
				trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(service, c.Path(), request)...),
				trace.WithSpanKind(trace.SpanKindServer),
			}

			ctx, span := tracing.DefaultTracer().Start(ctx, fmt.Sprintf("%s %s", request.Method, c.Path()), opts...)
			defer span.End()

			c.SetRequest(request.WithContext(ctx))
			c.Response().Header().Set("trace_id", span.SpanContext().TraceID().String())

			err := next(c)
			if err != nil {
				span.SetAttributes(attribute.String("echo.error", err.Error()))
				c.Error(err)
			}

			span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(c.Response().Status)...)
			span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(c.Response().Status, trace.SpanKindServer))

			return err
		}
	}
}
