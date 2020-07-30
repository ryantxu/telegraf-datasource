# Telegraf Datasource

Stream data directly from telegraph to grafana


1. Enable streaming in grafana.  
Run a version of grafana 7.2+ after July 27, 2020

In `conf/custom.ini`
```
[feature_toggles]
enable = live
```

2. Add telegraph output connector pointing to grafana server:
```
url: ws://localhost:3000/live/ws?format=protobuf
channel: telegraf
```

You will need to [build telegraf](./telegraf/README.md) with the [grafana-live](./telegraf/plugins/outputs/grafana-live/README.md) plugin. 
This can also be done from the root directory using
```
yarn build-telegraf
```

3. Add the datasource and configure with same channel configured above

4. Add a query pointing to the the right measuremnt
