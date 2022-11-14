# Cloudflare Metrics

## Development

Copy `local.env.dist` to `local.env` and set variables.

Use the following make commands:

```
make build
make lint
make test
```

## Releasing

```
goreleaser build --snapshot --rm-dist
```
