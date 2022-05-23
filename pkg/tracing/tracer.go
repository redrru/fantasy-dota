package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const name = "DefaultTracer"

func DefaultTracer() trace.Tracer {
	return otel.Tracer(name)
}
