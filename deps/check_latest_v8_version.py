#!/usr/bin/python3

# Script to compare current v8 version installed and the latest stable v8 version
# Returns 1 if the versions are the same
# Returns 0 if the latest stable version is higher than the current installed one

import requests
import sys

# api-endpoint
URL = "https://omahaproxy.appspot.com/all.json?os=linux&channel=stable"

def compare_v8_version(current_version, new_version):
  if current_version == new_version:
    return 1
  else:
    return 0

# Current version
v8_version_file = open("v8_version", "r")
current_v8_version_installed = v8_version_file.read().strip()

# Get latest version
r = requests.get(url = URL)

data = r.json()

latest_stable_v8_version = data[0]["versions"][0]["v8_version"].strip()

result = compare_v8_version(current_v8_version_installed, latest_stable_v8_version)
sys.exit(result)
