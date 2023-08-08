package tracing

import (
	"context"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
)

// Span is a wrapper that generates both Opentracing traces and Prometheus metrics.
type Span struct {
	start            time.Time
	span             opentracing.Span
	labels           []string
	requestCounter   *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec
}

func StartSpan(ctx context.Context, labels []string, requestCounter *prometheus.CounterVec,
	latencyHistogram *prometheus.HistogramVec) *Span {
	span := Span{
		start:            time.Now(),
		labels:           labels,
		requestCounter:   requestCounter,
		latencyHistogram: latencyHistogram,
	}
	span.span, _ = opentracing.StartSpanFromContext(ctx, strings.Join(labels, "_"))
	return &span
}

func (span *Span) SetTag(key, value string) {
	span.span.SetTag(key, value)
}

// Finish tarminates the span and observes metrics. Returns elapsed time in seconds.
func (span *Span) Finish() float64 {
	span.span.Finish()
	span.requestCounter.WithLabelValues(span.labels...).Inc()
	elapsed := time.Now().Sub(span.start).Seconds()
	span.latencyHistogram.WithLabelValues(span.labels...).Observe(elapsed)
	return elapsed
}
