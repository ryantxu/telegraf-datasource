package grafanalive

import (
	"github.com/grafana/grafana-plugin-sdk-go/live"
)

type measurementsCollector struct {
	ch    *live.GrafanaLiveChannel
	index []int
}
