package all

import (
	_ "github.com/influxdata/telegraf/plugins/outputs/http"
	_ "github.com/influxdata/telegraf/plugins/outputs/influxdb"
	_ "github.com/ryantxu/telegraf-datasource/telegraf/plugins/outputs/grafana-live"
)
