package grafanalive

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/influxdata/telegraf/plugins/serializers/json"
)

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
		// Set a default serializer. You should not use anything else.
		serializer, err := json.NewSerializer(time.Duration(1) * time.Millisecond)
		if err != nil {
			log.Fatal("Could not initialize a json serializer")
		}
		return &GrafanaLive{
			serializer: serializer,
		}
	})
}
