[![Live](https://img.shields.io/badge/live-entrvm.xyz-brightgreen)](https://entrvm.xyz)
[![Go](https://img.shields.io/badge/go-1.26.3-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-lightgrey)](LICENSE)

# do_1_server

A lightweight HTTP server written in Go, deployed on a DigitalOcean Ubuntu VM and served publicly at [https://entrvm.xyz](https://entrvm.xyz).

Built as a hands-on exercise in Go server fundamentals, Linux systems administration, and production deployment patterns.

---

## What it does

Serves HTTP endpoints over TLS. Currently exposes a single `/health` endpoint. The server reads its port from an environment variable, wraps all requests in a logging middleware, and encodes responses as JSON.

---

## Stack

| Layer | Tool |
|---|---|
| Language | Go |
| Process manager | systemd |
| Reverse proxy | nginx |
| TLS | Let's Encrypt via certbot |
| Host | DigitalOcean — Ubuntu 24.04 |
| Domain | entrvm.xyz |

---

## Project structure

```
do_1_server/
├── main.go              # Entry point — reads PORT, registers routes, starts server
├── go.mod               # Module definition
└── internal/
    ├── handler.go       # HTTP handlers
    ├── handler_test.go  # Handler tests
    └── logger.go        # Request logging middleware
```

### `main.go`

Entry point. Reads the `PORT` environment variable (defaults to `8080`), registers routes on a `ServeMux`, wraps the mux with the logger middleware, and starts the server.

### `internal/handler.go`

Defines HTTP handlers. Each handler is responsible for a single route. Currently implements:

- `HealthHandler` — responds to `GET /health` with `{"status":"ok"}`

### `internal/logger.go`

Request logging middleware. Wraps any `http.Handler`, records the start time, calls the next handler, then logs the method, path, and elapsed time to stderr via `log.Printf`. systemd captures this output through `journalctl`.

---

## Endpoints

| Method | Path | Response |
|---|---|---|
| GET | `/health` | `{"status":"ok"}` |

---

## Running locally

```bash
git clone https://github.com/bootupAbdullah/do_1_server
cd do_1_server
go build ./...
PORT=8080 ./do_1_server
```

Test:

```bash
curl localhost:8080/health
```

---

## Running tests

```bash
go test ./...
```

---

## Infrastructure

The server runs on a DigitalOcean Ubuntu 24.04 VM. Three OS-level components support it in production:

### systemd

Manages the server process. The binary is installed at `/usr/local/bin/do_1_server` and registered as a systemd service at `/etc/systemd/system/do_1_server.service`. systemd starts it on boot, restarts it on failure, and captures logs through `journalctl`.

```bash
systemctl status do_1_server
journalctl -u do_1_server -f
```

### nginx

Acts as a reverse proxy. nginx listens on ports 80 and 443, handles TLS termination, and forwards requests to the Go server on `localhost:8080`. The server is never exposed directly to the internet.

```
internet → nginx (443) → go server (8080)
```

Config lives at `/etc/nginx/sites-available/do_1_server`.

### certbot

Manages the TLS certificate from Let's Encrypt. Certbot modifies the nginx config automatically and installs a systemd timer to handle certificate renewal before expiry.

```bash
certbot renew --dry-run    # verify auto-renewal is working
certbot certificates       # list certificates and expiry dates
```

---

## Deployment

After making changes:

```bash
# on the VM
cd ~/do_1_server
git pull
go build ./...
cp do_1_server /usr/local/bin/do_1_server
systemctl restart do_1_server
```
