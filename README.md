# prometheus-solaxcloud-exporter

Prometheus exporter for [SolaxCloud](https://www.solaxcloud.com) data

## Install

Using Go 1.18 or newer

```shell
go install github.com/loafoe/prometheus-solaxcloud-exporter@latest
```

## Usage

### Set credentials

```shell
export SOLAXCLOUD_SN=your-sn
export SOLAXCLOUD_TOKEN_ID=your-token
```

### Run exporter

```shell
prometheus-solaxcloud-exporter -listen 0.0.0.0:8887
```

### Ship to prometheus

You can use something like Grafana-agent to ship data to a remote write endpoint. Example:

```yml
metrics:
  configs:
    - name: default
      scrape_configs:
        - job_name: 'solaxcloud_exporter'
          scrape_interval: 2m
          static_configs:
            - targets: ['localhost:8887']
      remote_write:
        - url: https://prometheus.example.com/api/v1/write
          basic_auth:
            username: scraper
            password: S0m3pAssW0rdH3Re
```

## License

License is MIT
