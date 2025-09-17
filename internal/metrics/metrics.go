package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    HTTPRequestDurationSeconds = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )

    WSConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "ws_connections",
            Help: "Number of active WebSocket connections",
        },
    )

    MessagesBroadcastTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "messages_broadcast_total",
            Help: "Total messages broadcasted per room",
        },
        []string{"room"},
    )
)


