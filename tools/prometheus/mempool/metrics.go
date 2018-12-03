package mempool

//import (
//	"github.com/cosmos/cosmos-sdk/client/context"
//	"github.com/tendermint/tendermint/mempool"
//
//	"log"
//	"time"
//	"strings"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//)
//
//// Metrics contains metrics exposed by this package.
//// see MetricsProvider for descriptions.
//type Metrics struct {
//	TmMetrics mempool.Metrics
//}
//
//// PrometheusMetrics returns Metrics build using Prometheus client library.
//func PrometheusMetrics() *Metrics {
//	tmMetrics := *mempool.PrometheusMetrics("tendermint")
//	return &Metrics{
//		tmMetrics,
//	}
//}
//
//func (m *Metrics) Start(rpc context.CLIContext) {
//	go func() {
//		for {
//			time.Sleep(1 * time.Second)
//			if result, err := rpc.(); err == nil {
//				m.TmMetrics.Size.Set(float64(result.N))
//			} else {
//				log.Println(err)
//			}
//		}
//	}()
//}
