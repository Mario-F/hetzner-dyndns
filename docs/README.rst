# Hetzner DynDNS

A script to update your dns entry on hetzner with your actual external IP.

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

.. _readme:

Hetzner DynDNS
================

|img_travis| |img_sr|

.. |img_travis| image:: https://travis-ci.org/Mario-F/hetzner-dyndns.svg?branch=master
   :alt: Travis CI Build Status
   :scale: 100%
   :target: https://travis-ci.com/Mario-F/promtail-formula
.. |img_sr| image:: https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg
   :alt: Semantic Release
   :scale: 100%
   :target: https://github.com/semantic-release/semantic-release

A tool to update Hetzner DNS with your external IP Address.

.. contents:: **Table of Contents**
   :depth: 1

Quickstart
----------


