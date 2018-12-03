package prometheus

import (
	cs 	"github.com/kidinamoto01/CosmosMonitor/tools/prometheus/consensus"
	sys "github.com/kidinamoto01/CosmosMonitor/tools/prometheus/system"
	gov "github.com/kidinamoto01/CosmosMonitor/tools/prometheus/governance"
	"github.com/kidinamoto01/CosmosMonitor/tools/prometheus/p2p"
	"github.com/cosmos/cosmos-sdk/client/context"
)

type Monitor struct {
	providers []MetricsProvider
	ctx       context.CLIContext
}

func DefaultMonitor(ctx context.CLIContext) *Monitor {
	var providers []MetricsProvider
	monitor := &Monitor{
		providers: providers,
		ctx:       ctx,
	}
	monitor.AddMetricsProvider(cs.PrometheusMetrics()).
		AddMetricsProvider(p2p.PrometheusMetrics()).
		//AddMetricsProvider(mempl.PrometheusMetrics()).
		AddMetricsProvider(sys.PrometheusMetrics()).
		AddMetricsProvider(gov.PrometheusMetrics())
	return monitor
}

func (m *Monitor) AddMetricsProvider(provider MetricsProvider) *Monitor {
	m.providers = append(m.providers, provider)
	return m
}

func (m *Monitor) Start() {
	for _, provider := range m.providers {
		provider.Start(m.ctx)
	}
}

type MetricsProvider interface {
	Start(ctx context.CLIContext)
}
