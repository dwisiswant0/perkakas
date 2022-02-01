package tracer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var DhuwitTracer *tracer

var (
	uintType = map[reflect.Kind]bool{
		reflect.Uint: true, reflect.Uint8: true, reflect.Uint16: true, reflect.Uint32: true, reflect.Uint64: true,
	}
	intType = map[reflect.Kind]bool{
		reflect.Int: true, reflect.Int8: true, reflect.Int16: true, reflect.Int32: true, reflect.Int64: true,
	}
	floatType = map[reflect.Kind]bool{
		reflect.Float32: true, reflect.Float64: true,
	}
	stringType = map[reflect.Kind]bool{
		reflect.String: true,
	}
	boolType = map[reflect.Kind]bool{
		reflect.Bool: true,
	}
	// TODO: implement type data mapping from interface to map slice
	// sliceType = map[reflect.Kind]bool{
	// 	reflect.Slice: true,
	// }
	// mapType = map[reflect.Kind]bool{
	// 	reflect.Map: true,
	// }
)

type Tracer interface {
	Trace() Span
}

type TraceProvider struct {
	JaegerProvider *tracesdk.TracerProvider
}

type Span interface {
	Start(ctx context.Context, opName string) (context.Context, Span)
	Finish(addTags ...map[string]interface{})
	Context() context.Context
	AddError(err error)
	AddAttr(k string, v interface{})
	AddLog(logName string, opts map[string]interface{})
	InjectHTTPClientHeader(ctx context.Context, header http.Header)
	ExtractHTTPClientHeader(ctx context.Context, header http.Header) context.Context
}

type Service struct {
	Name        string
	Version     string
	Environment string
}

type Option struct {
	Service           Service
	Enable            bool
	TraceProviderURL  string
	TraceProviderPort string
}

type tracer struct {
	Option        Option
	Otel          trace.Tracer
	traceProvider TraceProvider
}

type span struct {
	tracer    tracer
	tcSpan    trace.Span
	tcContext context.Context
}

func New(op Option) Tracer {
	tc := &tracer{
		Option: op,
	}

	DhuwitTracer = tc

	if !op.Enable {
		return tc
	}

	// tracer provider
	provider := TraceProvider{}
	jg, err := tc.JaegerProvider(op.TraceProviderURL, op.TraceProviderPort)
	if err != nil {
		log.Fatal(err)
	}

	// set provider
	provider.JaegerProvider = jg
	otel.SetTracerProvider(jg)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	))
	tc.Otel = otel.Tracer(op.Service.Name)

	return tc
}

func (tc *tracer) JaegerProvider(url string, port string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(url), jaeger.WithAgentPort(port)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(tc.Option.Service.Name),
			attribute.String("environment", tc.Option.Service.Environment),
			attribute.String("version", tc.Option.Service.Version),
		)),
		// TODO riset lebih untuk sampler
		// tracesdk.WithSampler(tracesdk.ParentBased())
	)
	return tp, nil
}

func (t *tracer) Trace() Span {
	return &span{
		tracer: *t,
	}
}

func (s *span) Start(ctx context.Context, opName string) (context.Context, Span) {
	if s.tracer.Option.Enable {
		pc, _, _, _ := runtime.Caller(1)
		fnName := filepath.Base(runtime.FuncForPC(pc).Name())
		requestId := ctx.Value(constants.CONTEXT_KEY_MSG_ID)
		s.tcContext, s.tcSpan = s.tracer.Otel.Start(ctx, opName, trace.WithAttributes(
			attribute.String("code.function", fnName),
			attribute.String("request.id", fmt.Sprintf("%v", requestId)),
		))
	} else {
		return ctx, s
	}
	return s.Context(), s
}

func (s *span) Finish(addAttrs ...map[string]interface{}) {
	if s.tcSpan == nil {
		return
	}

	defer s.tcSpan.End()
}

func (s *span) Context() context.Context {
	return s.tcContext
}

func (s *span) AddAttr(k string, v interface{}) {
	if s.tcSpan == nil {
		return
	}

	s.tcSpan.SetAttributes(s.setAttr(k, v))
}

func (s *span) setAttr(k string, v interface{}) attribute.KeyValue {
	kind := reflect.TypeOf(v).Kind()
	switch {
	case stringType[kind]:
		return attribute.String(k, v.(string))
	case floatType[kind]:
		return attribute.Float64(k, v.(float64))
	case intType[kind]:
		return attribute.Int64(k, v.(int64))
	case uintType[kind]:
		return attribute.Int64(k, int64(v.(uint64)))
	case boolType[kind]:
		return attribute.Bool(k, v.(bool))
	}

	return attribute.String(k, fmt.Sprint(v))
}

func (s *span) AddError(err error) {
	if s.tcSpan == nil {
		return
	}
	s.tcSpan.RecordError(err)
	s.tcSpan.SetStatus(codes.Error, err.Error())
}

func (s *span) AddLog(logName string, opts map[string]interface{}) {
	var kvs []attribute.KeyValue
	for k, v := range opts {
		kvs = append(kvs, s.setAttr(k, v))
	}
	s.tcSpan.AddEvent(logName, trace.WithAttributes(kvs...))
}

func (s *span) InjectHTTPClientHeader(ctx context.Context, header http.Header) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(header))
}

func (s *span) ExtractHTTPClientHeader(ctx context.Context, header http.Header) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(header))
}
