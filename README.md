# cofe

## Quick-start
Download and install:
```bash
curl -L https://github.com/mvsrgc/cofe/releases/download/v1.0/cofe -o cofe
chmod +x cofe
sudo mv cofe /usr/local/bin
```

Run a 4 minute timer (default):
`go run main.go`

Run a custom time :
`go run main.go 1h`
`go run main.go 15ms`
`go run main.go 0.5m`
`go run main.go 30s`
`go run main.go 15ms`
