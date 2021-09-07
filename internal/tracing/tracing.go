package tracing

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	traceFmt "github.com/opentracing/opentracing-go/log"
)

const parentSpan = "ParentSpan"

type Span struct {
	impl opentracing.Span
}

func RegisterSpan(ctx context.Context, name string) (context.Context, Span) {
	if parentSpan := ctx.Value(parentSpan); parentSpan != nil {
		if parent, success := parentSpan.(Span); success {
			return ctx, Span{
				impl: opentracing.StartSpan(name, opentracing.ChildOf(parent.impl.Context())),
			}
		}
	}
	span := opentracing.StartSpan(name)
	ctx = context.WithValue(ctx, parentSpan, span)
	return ctx, Span{impl: span}
}

func (s *Span) SpawnChild(parent Span, name string) Span {
	return Span{
		impl: opentracing.StartSpan(name, opentracing.ChildOf(parent.impl.Context())),
	}
}

func (s *Span) Finish() {
	if s.impl != nil {
		s.impl.Finish()
	}
}

func (s *Span) WriteError(err error) {
	s.impl.LogFields(traceFmt.String("error", err.Error()))
}

func (s *Span) WriteInfo(msg string, args ...interface{}) {
	value := fmt.Sprintf(msg, args...)
	s.impl.LogFields(traceFmt.String("info", value))
}
