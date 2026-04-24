# logpipe

A lightweight CLI tool for tailing, filtering, and forwarding structured logs from multiple sources in real time.

---

## Installation

```bash
go install github.com/yourusername/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logpipe.git
cd logpipe
go build -o logpipe .
```

---

## Usage

Tail a log file and filter by log level:

```bash
logpipe tail --file /var/log/app.log --filter level=error
```

Forward logs from multiple sources to a remote endpoint:

```bash
logpipe forward \
  --source /var/log/app.log \
  --source /var/log/worker.log \
  --output http://logs.example.com/ingest \
  --format json
```

Pipe stdin directly through logpipe:

```bash
journalctl -f | logpipe filter --level warn
```

### Common Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to a log file to tail |
| `--filter` | Key=value filter expression |
| `--format` | Output format: `json`, `text` (default: `json`) |
| `--output` | Forward logs to a URL or file path |
| `--level` | Minimum log level to display |

---

## License

MIT © 2024 yourusername