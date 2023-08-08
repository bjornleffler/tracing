package tracing

// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	parent           *Span
	requestCounter   *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec
	// TODO(leffler): Server exclusive latency.
}

func StartServerSpan(ctx context.Context, labels []string) *Span {
	span := Span{
		start:            time.Now(),
		labels:           labels,
		parent:           nil,
		requestCounter:   serverRequests,
		latencyHistogram: serverLatency,
	}
	span.span, _ = opentracing.StartSpanFromContext(ctx, strings.Join(labels, "_"))
	return &span
}

func StartClientSpan(ctx context.Context, parent *Span, labels []string) *Span {
	span := Span{
		start:            time.Now(),
		labels:           labels,
		parent:           parent,
		requestCounter:   clientRequests,
		latencyHistogram: clientLatency,
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
	// TODO(leffler): Update parent span.
	return elapsed
}
