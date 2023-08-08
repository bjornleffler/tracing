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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	labels         = []string{"service", "method"}
	clientRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "client_request_total",
		Help: "The total number of client requests, by service and method."},
		labels)
	clientLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "client_latency",
		Help: "Client request latency, by service and method."},
		labels)

	serverRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "server_request_total",
		Help: "The total number of client requests, by service and method."},
		labels)
	serverLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "server_latency",
		Help: "Server request latency, by service and method."},
		labels)

	// "Exclusive" latency is the strict server latency, excluding client latency.
	serverExclusiveLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "server_exclusive_latency",
		Help: "Server exclusive request latency, by service and method."},
		labels)
)
