# Grafana Live streaming output plugin
## Configuration
Example
```
[[outputs.grafana-live]]
  # This points to your running grafana instance
  address = "localhost:3000"
  # The channel that you want to publish messages to
  channel = "telegraf"
```
