# cofe

## Preview
[![asciicast](https://asciinema.org/a/5iMyxSZoAIYGRFms4a7FXbx8J.svg)](https://asciinema.org/a/5iMyxSZoAIYGRFms4a7FXbx8J)

## Quick-start
Download and install:
```bash
curl -L https://github.com/mvsrgc/cofe/releases/download/v1.0/cofe -o cofe
chmod +x cofe
sudo mv cofe /usr/local/bin
```

Run a 4 minute timer (default):
`cofe`

Run a custom time:
`cofe 1h`
`cofe 15ms`
`cofe 0.5m`
`cofe 30s`
`cofe 15ms`

Output to stdout with no decorations (and write to /tmp/cofe_status):
`cofe --raw 1m`
