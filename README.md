# portwatch

Lightweight daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourname/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/portwatch.git && cd portwatch && go build ./...
```

## Usage

Start the daemon with a config file:

```bash
portwatch --config /etc/portwatch/config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
baseline:
  - 22
  - 80
  - 443
alerts:
  slack_webhook: "https://hooks.slack.com/services/..."
```

On first run, `portwatch` snapshots the currently open ports as the baseline. Any port that opens or closes unexpectedly will trigger an alert.

```bash
# Run once and print current open ports
portwatch --scan-once

# Run in foreground with verbose logging
portwatch --config config.yaml --verbose
```

## How It Works

1. Polls open TCP/UDP ports at a configurable interval
2. Compares results against the known baseline
3. Fires an alert (log, webhook, or email) when a diff is detected

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

## License

MIT © yourname