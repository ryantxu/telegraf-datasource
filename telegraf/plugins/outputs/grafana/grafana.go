package grafana

import (
	"fmt"
	"log"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/influxdata/telegraf/plugins/serializers/json"
	"github.com/ryantxu/telegraf-datasource/telegraf/tds"
)

type Grafana struct {
	Address string `toml:"address"`
	Channel string `toml:"channel"`

	broker     *tds.GrafanaLiveChannel
	serializer serializers.Serializer
}

var sampleConfig = `
[[outputs.grafana]]
  # The address of the local grafana instance
  address = "localhost:3000"
  # The channel to write data into grafana with
  channel = "telegraf"
`

func (g *Grafana) Connect() error {
	var err error
	g.broker, err = tds.InitGrafanaLiveChannel(fmt.Sprintf("ws://%s/live/ws?format=protobuf", g.Address, g.Format), g.Channel)
	if err != nil {
		return err
	}

	return err
}

func (g *Grafana) Close() error {
	return nil
}

func (g *Grafana) SampleConfig() string {
	return sampleConfig
}

func (g *Grafana) Description() string {
	return "Send telegraf metrics to a grafana live stream"
}

func (g *Grafana) Write(metrics []telegraf.Metric) error {
	b, err := g.serializer.SerializeBatch(metrics)
	if err != nil {
		return err
	}

	g.broker.Publish(b)

	return nil
}

func init() {
	outputs.Add("grafana", func() telegraf.Output {
		// Set a default serializer. You should not use anything else.
		serializer, err := json.NewSerializer(time.Duration(1) * time.Millisecond)
		if err != nil {
			log.Fatal("Could not initialize a json serializer")
		}
		return &Grafana{
			serializer: serializer,
		}
	})
}
