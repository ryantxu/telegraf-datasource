# Telegraf build + plugins

In order to stream data to grafana from telegraf you will need a build of telegraf that contains the grafana-live plugin, which is included in this repo

## Including plugins for build

The build (currently only supported on linux/mac) will clone the telegraf repository, and massage the configuration in this repository. 
If you want to add a plugin, you need to add it to the plugin type under `all/extra.go`

```
package all

import (
	_ "github.com/ryantxu/telegraf-datasource/telegraf/plugins/inputs/replay"
	_ "github.com/srclosson/telegraf-plugins/plugins/inputs/dirmon"
	_ "github.com/srclosson/telegraf-plugins/plugins/inputs/nmea"
)
```

The file above has the default import for the replay plugin, but also imports two other plugins called "dirmon" and "nmea".

Once the `extra.go` file has been updated, you need to add an entry to `etc/replacement.go.mod` to include the github repo where the module is located

```
replace (
	github.com/ryantxu/telegraf-datasource => /Users/stephanie/src/plugins/telegraf-datasource
	github.com/influxdata/telegraf => /tmp/telegraf
	github.com/srclosson/telegraf-plugins => /tmp/telegraf-plugins
	github.com/satori/go.uuid => github.com/gofrs/uuid v3.2.0+incompatible
)
```

The line added is: `github.com/srclosson/telegraf-plugins => /tmp/telegraf-plugins`. 
Reference the github repo to include in the first section, and in the section behind the `=>` reference the location you want to clone the repo, so the build can find it.

## Building
### Linux
```
make
```

or from the root directory of this project
```
yarn build-telegraf
```

### MacOS
The same build steps as for Linux apply to MacOS with an additional option. To build for Linux on MacOS
```
make docker-linux
```

This will download a docker build image, and execute the build within that docker image.

## Troubleshooting
Often older copies of cached libraries can be a problem. Try cleaning your GOPATH
```
/bin/rm -rf $(go env GOPATH)
```

