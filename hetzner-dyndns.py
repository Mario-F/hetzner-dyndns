#!/usr/bin/env python
#
# Update a DNS Record on Hetzner DNS API with your external IP Address
#
# Install Python libraries:
# `pip install -r requirements.txt`
#

import argparse, requests, re, time, random, json

version = 0.1

# Parse commandline arguments
parser = argparse.ArgumentParser(description="Hetzner Custom dyndns updates v" + str(version))
parser.add_argument("mode", help="Mode to use: update (update dns), records (get all zones and list records)")
parser.add_argument("-t", "--token", help="Your Hetzner DNS API-Token", required=True)
parser.add_argument("-r", "--record", help="Record ID to update")
parser.add_argument("-q", "--quite", help="Dont print info messages", action='store_true')
parser.add_argument("-d", "--debug", help="Show debug information", action='store_true')
args = parser.parse_args()

#
# Main Program
#
def main():
  # Mode - Records
  if args.mode == 'records':
    print('Get and list all records from all zones:')
    for record in _getAllRecords():
      print("ID: {} - {}.{}".format(record['id'], record['name'], record['zone']['name']))
    exit(0)

  # Mode - Update
  if args.mode == 'update':
    # Get actual record with value
    actualRecord = _getRecord(args.record)
    _log('Actual value for domain {}: {}'.format(actualRecord['name'], actualRecord['value']))

    # Randomized list of IP Provider Functions
    useIpProviders = [_ipIpify, _ipCheckIP, _ipIfconfigME]
    random.shuffle(useIpProviders)

    retIPs = []
    foundIP = None
    for provider in useIpProviders:
      try:
        ip = provider()
        # Check if IP already returned from other provider
        # If so the IP should be valid
        if ip in retIPs:
          foundIP = ip
          _log(f"Valid IP found: {ip}")
          break
        else:
          retIPs.append(ip)
      except:
        _log('An error occured')

    # No ip found
    if not foundIP:
      _log('There was a problem lookup outside IP.')
      exit(1)
    
    # If IP unchanged exit sucessfully
    if foundIP == actualRecord['value']:
      _log('IP not changed, exit.')
      exit(0)

    # Update IP Address
    actualRecord['value'] = foundIP
    if _putRecord(actualRecord):
      _log('IP updated successfully')
      exit(0)

    _log('There was an error updating the IP address')
    exit(1)


  # Mode not Found
  print('Please provide an valid mode!')
  exit(1)


#
# Hetzner API
#
def _apiGet(endpoint):
  token = args.token
  response = requests.get(url=f"https://dns.hetzner.com/api/v1/{endpoint}", headers={"Auth-API-Token": token})
  _debug('Response HTTP Status Code: {status_code}'.format(status_code=response.status_code))
  _debug('Response HTTP Response Body: {content}'.format(content=response.content))
  return response.json()

def _getAllRecords():
  response = _apiGet('zones')['zones']
  zones = { z['id'] : z for z in response }
  records = list(filter(lambda x: x['type'] == 'A', _apiGet('records')['records']))
  for record in records:
    record['zone'] = zones[record['zone_id']]
  return records

def _getRecord(record):
  record = _apiGet(f'records/{record}')['record']
  return record

def _putRecord(record):
  token = args.token
  updateRecord = {
    "value": record['value'],
    "type": record['type'],
    "name": record['name'],
    "zone_id": record['zone_id']
  }
  if record.get('ttl'):
    updateRecord['ttl'] = record.get('ttl')
  response = requests.put(
    url="https://dns.hetzner.com/api/v1/records/{}".format(record['id']), 
    headers={"Content-Type": "application/json", "Auth-API-Token": token},
    data=json.dumps(updateRecord)
  )
  print('Response HTTP Status Code: {status_code}'.format(status_code=response.status_code))
  print('Response HTTP Response Body: {content}'.format(content=response.content))
  if response.status_code == 200:
    return True
  return False

#
# Helpers
#
def _log(msg):
  if args.quite != True:
    print('LOG - ' + time.strftime("%H:%M:%S") + ': ' + str(msg))

def _debug(msg):
  if args.debug == True:
    print('DEBUG - ' + time.strftime("%H:%M:%S") + ': ' + str(msg))

def _getJSON(url):
  r = requests.get(url)
  return r.json()

def _captureIP(ip):
  res = re.search('(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)', ip)
  if res is None:
    return None
  return str(res.group(0))


#
# IP Providers
#
def _ipIpify():
  res = _getJSON('https://api.ipify.org?format=json')
  ip = _captureIP(res['ip'])
  if not ip:
    _log('Failed to aquire IP from ipify')
    return False
  _log(f'IP from ipify: {ip}')
  return ip

def _ipCheckIP():
  res = requests.get('http://checkip.dyndns.org')
  ip = _captureIP(res.text)
  if not ip:
    _log('Failed to aquire IP from checkip')
    return False
  _log(f'IP from checkip: {ip}')
  return ip

def _ipIfconfigME():
  res = requests.get('http://ifconfig.me')
  ip = _captureIP(res.text)
  if not ip:
    _log('Failed to aquire IP from ifconfigME')
    return False
  _log(f'IP from ifconfigME: {ip}')
  return ip

if __name__ == "__main__":
    main()