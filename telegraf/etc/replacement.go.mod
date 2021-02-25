module github.com/influxdata/telegraf

go 1.14

// replaced due to https://github.com/satori/go.uuid/issues/73
replace (
	github.com/ryantxu/telegraf-datasource => /home/ryan/workspace/grafana/more/telegraf-datasource
	github.com/influxdata/telegraf => /tmp/telegraf
	github.com/satori/go.uuid => github.com/gofrs/uuid v3.2.0+incompatible
)


// Trial and error replaced values with values from go.mod telegraf master
require (
	github.com/influxdata/toml v0.0.0-20190415235208-270119a8ce65
	github.com/soniah/gosnmp v1.25.0
	github.com/gosnmp/gosnmp v1.29.0
)
