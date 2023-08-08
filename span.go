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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Span is a wrapper that generates both Opentracing traces and Prometheus metrics.
type Span struct {
	start            time.Time
	span             opentracing.Span
	service          string
	method           string
	parent           *Span
	requestCounter   *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec
	clientElapsed    float64
}

var defaultService = ""

func Configure(service string, port int) {
	defaultService = service
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func StartServerSpan(ctx context.Context, method string) *Span {
	span := Span{
		start:            time.Now(),
		service:          defaultService,
		method:           method,
		parent:           nil,
		requestCounter:   serverRequests,
		latencyHistogram: serverLatency,
	}
	span.span, _ = opentracing.StartSpanFromContext(ctx, strings.Join(labels, "_"))
	return &span
}

func StartClientSpan(ctx context.Context, parent *Span, service, method string) *Span {
	span := Span{
		start:            time.Now(),
		service:          service,
		method:           method,
		parent:           parent,
		requestCounter:   clientRequests,
		latencyHistogram: clientLatency,
		clientElapsed:    0,
	}
	span.SetTag("span.kind", "client")
	span.span, _ = opentracing.StartSpanFromContext(ctx, strings.Join(labels, "_"))
	return &span
}

func (span *Span) SetTag(key, value string) {
	span.span.SetTag(key, value)
}

// Finish terminates the span and observes metrics. Returns elapsed time in seconds.
func (span *Span) Finish() float64 {
	span.span.Finish()
	span.requestCounter.WithLabelValues(span.service, span.method).Inc()
	elapsed := time.Now().Sub(span.start).Seconds()
	span.latencyHistogram.WithLabelValues(span.service, span.method).Observe(elapsed)
	if span.parent == nil {
		exclusive := elapsed - span.clientElapsed
		serverExclusiveLatency.WithLabelValues(span.service, span.method).Observe(
			exclusive)
	} else {
		span.parent.clientElapsed += elapsed
	}
	return elapsed
}
