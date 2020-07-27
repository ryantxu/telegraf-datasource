package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ryantxu/telegraf-datasource/pkg/tds"
)

func main() {
	// Background thread
	fmt.Println("starting http server")

	b, err := tds.InitGrafanaLiveChannel("ws://localhost:3000/live/ws?format=protobuf", "sinal-over-ws")
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
	streamSignal(b)
}

// write to a stream....
func streamSignal(broker *tds.GrafanaLiveChannel) {
	speed := 1000 / 1 // 20 hz
	spread := 50.0

	walker := rand.Float64() * 100
	ticker := time.NewTicker(time.Duration(speed) * time.Millisecond)

	line := tds.InfluxLine{
		Name:   "simple",
		Fields: make(map[string]interface{}),
		Tags:   make(map[string]string),
	}

	s, _ := tds.NewSerializer(time.Duration(1) * time.Millisecond)

	for t := range ticker.C {
		delta := rand.Float64() - 0.5
		walker += delta

		line.Timestamp = t
		line.Fields["value"] = walker
		line.Fields["min"] = walker - ((rand.Float64() * spread) + 0.01)
		line.Fields["max"] = walker + ((rand.Float64() * spread) + 0.01)

		b, _ := s.SerializeBatch([]*tds.InfluxLine{&line}) //json.Marshal(line)

		broker.Publish(b)
	}
}
