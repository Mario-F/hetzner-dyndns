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
