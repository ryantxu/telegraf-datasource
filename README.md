# Telegraf Datasource

Stream data directly from telegraph to grafana


1. Enable streaming in grafana.  In `conf/custom.ini`
```
[feature_toggles]
enable = live
```

2. Add telegraph output connector pointing to grafana server:
```
url: ws://localhost:3000/live/ws?format=protobuf
channel: telegraf
```

3. Add the datasource and configure with same channel configured above

4. Add a query pointing to the the right measuremnt
