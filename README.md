# TCP(+HTTP) Reverse Proxy in GO
TCP reverse proxy written in go trying to keep the project using standard libarary as much as possible
Features:
- endpoint routing (HTTP mode)
- change header (HTTP mode)
- rate limit (HTTP mode) <- NOT DONE
- loadbalancing (Random, IP hashing, least-connection, round-robin)
- healthcheck (no self healing) 

TODO:
- Handle level 4 connection
- Rate limit for HTTP mode
- Refractor code
- Do logging properly
### Quick Start
```bash
git clone https://github.com/phoon884/rev-proxy
cp example_config.yaml config.yaml # Or write a new one
go run ./cmd/main.go
```
