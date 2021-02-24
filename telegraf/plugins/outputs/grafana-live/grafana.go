package grafanalive

import (
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/live"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
)

// GrafanaLive connects to grafana server
type GrafanaLive struct {
	Address string `toml:"address"`
	Channel string `toml:"channel"`

	client     *live.GrafanaLiveClient
	channels   map[string]*live.GrafanaLiveChannel
	serializer serializers.Serializer
}

var sampleConfig = `
[[outputs.grafana]]
  # The address of the local grafana instance
  address = "localhost:3000"
  # The channel to write data into grafana with
  channel = "telegraf"
`

func (g *GrafanaLive) Connect() error {
	var err error

	g.client, err = live.InitGrafanaLiveClient(live.ConnectionInfo{
		URL: g.Address,
	})
	if err != nil {
		return err
	}
	g.channels = make(map[string]*live.GrafanaLiveChannel)
	g.client.Subscribe(live.ChannelAddress{})

	return err
}

func (g *GrafanaLive) Close() error {

	return nil
}

func (g *GrafanaLive) SampleConfig() string {
	return sampleConfig
}

func (g *GrafanaLive) Description() string {
	return "Send telegraf metrics to a grafana live stream"
}

func (g *GrafanaLive) getChannel(name string) *live.GrafanaLiveChannel {
	c, ok := g.channels[name]
	if ok {
		return c
	}

	var err error
	c, err = g.client.Subscribe(live.ChannelAddress{
		Scope:     "grafana",
		Namespace: "measurments",
		Path:      g.Channel + "/" + name,
	})
	if err != nil {
		fmt.Println("ERROR")
	}
	return c
}

type measurementsCollector struct {
	ch      *live.GrafanaLiveChannel
	metrics []telegraf.Metric
}

func (g *GrafanaLive) Write(metrics []telegraf.Metric) error {
	measures := make(map[string]measurementsCollector)
	for _, metric := range metrics {
		name := metric.Name()
		m, ok := measures[name]
		if !ok {
			m = measurementsCollector{
				ch: g.getChannel(name),
			}
			measures[name] = m
		}
		m.metrics = append(m.metrics, metric)
	}

	for _, val := range measures {
		if val.ch != nil {

			b, err := g.serializer.SerializeBatch(val.metrics)
			if err != nil {
				return err
			}

			val.ch.Publish(b)
		}
	}

	return nil
}

func init() {
	outputs.Add("grafana-live", func() telegraf.Output {
		// This is the serializer that grafana will understand
		s := &serializer{
			TimestampUnits: int64(time.Duration(1) * time.Millisecond),
		}
		return &GrafanaLive{
			serializer: s,
		}
	})
}
