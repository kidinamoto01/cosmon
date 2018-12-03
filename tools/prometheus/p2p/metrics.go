package p2p

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/pelletier/go-toml"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
	"github.com/cosmos/cosmos-sdk/client/context"
	"fmt"
	"net/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

)

type Metrics struct {

	// Number of peers.
	Peers metrics.Gauge
	// Number of connected persistent peers.
	ConnectedPersistentPeers metrics.Gauge
	// Number of unconnected persistent peers.
	UnonnectedPersistentPeers metrics.Gauge
	persistent_peers          map[string]string
}

func PrometheusMetrics() *Metrics {
	return &Metrics{
		Peers: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "p2p",
			Name:      "peers",
			Help:      "Number of peers.",
		}, []string{}),
		ConnectedPersistentPeers: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "p2p",
			Name:      "connected_persistent_peers",
			Help:      "Number of connected persistent peers.",
		}, []string{}),
		UnonnectedPersistentPeers: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "p2p",
			Name:      "unconnected_persistent_peers",
			Help:      "Number of unconnected persistent peers.",
		}, []string{}),
		persistent_peers: make(map[string]string),
	}
}

func (m *Metrics) Start(ctx context.CLIContext) {
	m.setP2PPersistentPeers(viper.GetString("home"))
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if result, err := m.getNetInfo(ctx); err == nil {
				connected := 0
			for _, peer := range result.Peers {
					if _, exist := m.persistent_peers[string(peer.NodeInfo.ID_)]; exist {
						connected += 1
					}
				}
				m.Peers.Set(float64(result.NPeers))
				m.ConnectedPersistentPeers.Set(float64(connected))
				m.UnonnectedPersistentPeers.Set(float64(len(m.persistent_peers) - connected))
			} else {
				log.Println(err)
			}
		}
	}()
}

//set the p2p persistent peers by given home dir of iris config file
func (m *Metrics) setP2PPersistentPeers(homeDir string) {
	if !filepath.IsAbs(homeDir) {
		absHomeDir, err := filepath.Abs(homeDir)
		if err != nil {
			log.Println("cannot find the file ", err)
			return
		}
		homeDir = absHomeDir
	}
	configFilePath := filepath.Join(homeDir, "config/config.toml")
	//fmt.Printf("configFilePath: %s\n", configFilePath)
	if data, err := ioutil.ReadFile(configFilePath); err != nil {
		log.Println("cannot open the file ", err)
		return
	} else {
		if config, err := toml.LoadBytes(data); err != nil {
			log.Println("parse config file failed: ", err)
			return
		} else {
			persistent_peers := config.Get("p2p.persistent_peers").(string)
			for _, peer := range strings.Split(persistent_peers, ",") {
				if peer != "" {
					splited := strings.Split(peer, "@")
					m.persistent_peers[splited[0]] = splited[1]
				}
			}
		}
	}
}

func (m *Metrics) getNetInfo(ctx context.CLIContext)(*ctypes.ResultNetInfo, error){
	client := &http.Client{}


	url := strings.Replace(ctx.NodeURI, "tcp", "http", 1)
	reqUri := fmt.Sprintf("%s/%s", url, "num_unconfirmed_txs")

	resp, err := client.Get(reqUri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res = struct {
		JsonRpc string                      `json:"jsonrpc"`
		Id      string                      `json:"id"`
		Result  ctypes.ResultNetInfo `json:"result"`
	}{}


	if err := ctx.Codec.UnmarshalJSON(body, &res); err != nil {
		return nil, err
	}

	return &res.Result, nil
}