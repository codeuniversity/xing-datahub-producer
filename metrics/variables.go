package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

//HTTPProcessed counts the processed http requests
var HTTPProcessed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_processed",
		Help: "The amount of http requests processed partioned by status code and record type",
	},
	[]string{"code", "record"},
)
