//go:build unit
// +build unit

package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func testTraceProvider() (*trace.TracerProvider, error) {
	e, err := stdouttrace.New(
		stdouttrace.WithWriter(bytes.NewBuffer([]byte{})),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		return nil, err
	}

	p := trace.NewTracerProvider(
		trace.WithBatcher(e),
		trace.WithResource(resource.Default()),
	)

	return p, nil
}

func TestTracingFields(t *testing.T) {
	type args struct {
		ctx    context.Context
		fields []zap.Field
	}
	type want struct {
		result []zap.Field
	}

	fakeField := zap.Any(gofakeit.Word(), gofakeit.Int32())

	tp, err := testTraceProvider()
	assert.NoError(t, err)
	otel.SetTracerProvider(tp)
	traceCtx, span := otel.Tracer("test").Start(context.Background(), "Test")
	defer span.End()

	testCases := []struct {
		name string
		args args
		want want
	}{
		{
			name: "NilCtxNilFields",
			args: args{
				ctx:    nil,
				fields: nil,
			},
			want: want{
				result: []zap.Field{},
			},
		},
		{
			name: "NilCtx",
			args: args{
				ctx:    nil,
				fields: []zap.Field{fakeField},
			},
			want: want{
				result: []zap.Field{fakeField},
			},
		},
		{
			name: "CtxWithoutSpan",
			args: args{
				ctx:    context.Background(),
				fields: []zap.Field{fakeField},
			},
			want: want{
				result: []zap.Field{fakeField},
			},
		},
		{
			name: "CtxWithSpan",
			args: args{
				ctx:    traceCtx,
				fields: []zap.Field{fakeField},
			},
			want: want{
				result: []zap.Field{zap.String(traceID, span.SpanContext().TraceID().String()), zap.String(spanID, span.SpanContext().SpanID().String()), fakeField},
			},
		},
		{
			name: "CtxWithSpanNilFields",
			args: args{
				ctx:    traceCtx,
				fields: nil,
			},
			want: want{
				result: []zap.Field{zap.String(traceID, span.SpanContext().TraceID().String()), zap.String(spanID, span.SpanContext().SpanID().String())},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := withTracingFields(tc.args.ctx, tc.args.fields...)
			assert.Equal(t, tc.want.result, result)
		})
	}
}
