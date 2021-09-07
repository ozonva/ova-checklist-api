package application

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"
	jcfg "github.com/uber/jaeger-client-go/config"
	jlog "github.com/uber/jaeger-client-go/log"
	jmet "github.com/uber/jaeger-lib/metrics"

	"github.com/ozonva/ova-checklist-api/internal/config"
)

func startTracing(cfg *config.TraceConfig) io.Closer {
	jaegerConfig := jcfg.Configuration{
		ServiceName: cfg.ServiceName,
		Disabled:    !cfg.Enabled,
		Sampler: &jcfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jcfg.ReporterConfig{
			LogSpans: cfg.LogSpans,
		},
	}

	tracer, closer, err := jaegerConfig.NewTracer(
		jcfg.Logger(jlog.StdLogger),
		jcfg.Metrics(jmet.NullFactory),
	)

	if err != nil {
		log.Error().
			Str("reason", "cannot run the tracing subsystem").
			Msgf("%v", err)
		doCrash()
	}

	opentracing.SetGlobalTracer(tracer)

	return closer
}

func stopTracing(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Warn().
			Str("reason", "cannot stop tracing subsystem").
			Msgf("%v", err)
	}
}
