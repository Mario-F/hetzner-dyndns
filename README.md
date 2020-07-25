# Hetzner DynDNS

A script to update your dns entry on hetzner with your actual outside IP.

## IP Scraping

This script uses a random list of 'what-is-my-ip' providers to cut down on error.

## Usage

Clone this repo and install requirements with:

```shell
pip install -r requirements.txt
```

Get an DNS API-Token at Hetzner.

Query your Domains with:

```shell
./hetzner-dyndns.py --token 'your-api-token' records
```

Now execute again with mode update and the record ID for the domain to update with your outside IP:

```shell
./hetzner-dyndns.py --token 'your-api-token' --record 'domain-record-id' update
```
