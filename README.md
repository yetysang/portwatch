# portwatch

A lightweight CLI daemon that monitors port usage changes and alerts on unexpected bindings.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with default settings:

```bash
portwatch start
```

Watch specific ports and get alerted on unexpected bindings:

```bash
portwatch start --ports 80,443,8080 --interval 5s --alert email
```

Run a one-time snapshot of current port bindings:

```bash
portwatch scan
```

### Example Output

```
[2024-01-15 10:32:01] INFO  Watching 3 ports (interval: 5s)
[2024-01-15 10:32:06] WARN  Unexpected binding detected: 0.0.0.0:8443 (PID 3821 - unknown)
[2024-01-15 10:32:06] INFO  Port 443 still bound to nginx (PID 1042) — OK
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--ports` | all | Comma-separated list of ports to watch |
| `--interval` | `10s` | Polling interval |
| `--alert` | `stdout` | Alert method: `stdout`, `email`, `webhook` |
| `--config` | `~/.portwatch.yaml` | Path to config file |

---

## Configuration

You can define watched ports and alert settings in `~/.portwatch.yaml`:

```yaml
interval: 10s
ports:
  - 80
  - 443
  - 8080
alert: stdout
```

---

## License

MIT © [yourusername](https://github.com/yourusername)