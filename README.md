This example application generates OTel metrics in a simple client that get sent to the OTel Collector where they are exported to a file using the file exporter.

To run this example: `docker-compose up --build`

the hit the client endpoint: `curl localhost:8080`

The format of the output is specified in `otel-collector-config.yaml` and can be set to either `json` or `proto`. The output will be written to `file-exporter/metrics.<format>`.
