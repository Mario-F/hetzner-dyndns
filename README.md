# Hetzner DynDNS

A tool for updating Hetzner DNS with your external IP.

## IP Scraping

This program uses a random list of 'what-is-my-ip' providers to cut down on errors.

## Install

Install binaries for your system from the release page or use the install script:

```shell
# Local (current directory)
curl -s https://raw.githubusercontent.com/Mario-F/hetzner-dyndns/master/install.sh | BINDIR=./ bash

# Global
curl -s https://raw.githubusercontent.com/Mario-F/hetzner-dyndns/master/install.sh | sudo BINDIR=/usr/local/bin bash
```

## Usage

1. Test external IP with: `./hetzner-dyndns myip`
2. Obtain your API-Token from Hetzner DNS managment site.
3. Fetch your domain record ids with: `./hetzner-dyndns --token 'your-api-token' records`
4. Call update with record ID: `./hetzner-dyndns --token 'your-api-token' --record 'domain-record-id' update`

There is also a docker image available: `docker run --rm -it ghcr.io/mario-f/hetzner-dyndns:latest myip`

## Development / Debugging

There a serveral ways to easy start development and using live debugging provided by [delve](https://github.com/go-delve/delve)

### VSCode integrated Console

The provided launch.json has a debug task `Launch File` predefined, just hit start and it should run with the args provided in launch.json.

## VSCode external Terminal

A more advanced way to test in an external Terminal is provided by the `External Debugging` launch config and `./debug` script:

1. Execute debug script with optional arguments: `./debug --token super-secret-token records`
2. Start the `External Debugging` session in vscode

Unfortunately the order is importend because vscode does not try to automatically connect after start.
