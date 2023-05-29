package Jgrpc_otelspan

import (
	"context"
	"errors"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// OtelSpan 链路追踪
type OtelSpan struct {
	tracer *sdktrace.TracerProvider
}

// New 实例化 OtelSpan
func New(tracer *sdktrace.TracerProvider) *OtelSpan {
	return &OtelSpan{
		tracer: tracer,
	}
}

func (s *OtelSpan) Record(ctx context.Context, tracerName string) (context.Context, trace.Span) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return s.tracer.Tracer(tracerName).Start(ctx, f.Name())
}

// Error 记录链路错误
// 包含 文件名以及完整路径/错误行数
func (s *OtelSpan) Error(span trace.Span, msg string) error {
	_, file, line, _ := runtime.Caller(1)
	span.SetStatus(otelcodes.Error, msg)
	span.SetAttributes(
		attribute.String("file.name", file),
		attribute.Int("file.line", line),
	)
	return errors.New("msg")
}
