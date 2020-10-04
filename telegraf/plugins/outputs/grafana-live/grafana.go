package grafanalive

import (
	"fmt"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
)

// GrafanaLive connects to grafana server
type GrafanaLive struct {
	Address string `toml:"address"`
	Channel string `toml:"channel"`

	broker     *GrafanaLiveChannel
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
	g.broker, err = InitGrafanaLiveChannel(fmt.Sprintf("ws://%s/live/ws?format=protobuf", g.Address), g.Channel)
	if err != nil {
		return err
	}

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

func (g *GrafanaLive) Write(metrics []telegraf.Metric) error {
	b, err := g.serializer.SerializeBatch(metrics)
	if err != nil {
		return err
	}

	g.broker.Publish(b)

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
