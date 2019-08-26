package keeper

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	cfg "github.com/orientwalt/tendermint/config"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const MetricsSubsystem = "module_distribution"

type Metrics struct {
	CommunityTax metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}
	return &Metrics{
		CommunityTax: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: config.Namespace,
			Subsystem: MetricsSubsystem,
			Name:      "community_tax",
			Help:      "community tax",
		}, []string{}),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		CommunityTax: discard.NewGauge(),
	}
}
