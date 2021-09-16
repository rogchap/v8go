#!/bin/bash -e

echo "Cloning the v8 projects submodules..."
git submodule update --init --recursive

echo "Adding depot_tools folder to local PATH"
PATH=$PWD/deps/depot_tools:$PATH

echo "Finding the current stable release of v8..."
v8_versions=$(curl -sS 'https://omahaproxy.appspot.com/all.json?os=linux&channel=stable' | jq -r '.[0]["versions"][0]')
v8_commit="$(echo ${v8_versions} | jq -r .v8_commit)"
v8_version="$(echo ${v8_versions} | jq -r .v8_version)"

IFS='.'
read -ra ADDR <<< "$v8_version"
v8_major=${ADDR[0]}
v8_minor=${ADDR[1]}
IFS=' '

branch="v${v8_major}_${v8_minor}_upgrade"

echo "Creating upgrade branch ${branch}..."
git checkout -b ${branch}

pushd deps
  echo "Fetching the v8 source..."
  fetch v8

  pushd v8
    git checkout "${v8_commit}"
  popd


  echo "Copying contents of deps/v8/include to deps/include..."
  cp -r v8/include/* include/

  # For users that use go mod vendor, a vendor.go file is required in
  # subdirectories of deps/include. Do not remove any that exist, and
  # include it for any new directories.
  # TODO: Check git diff for new directories in deps/include to copy a vendor.go into
  # and update cgo.go

  echo "Building the binary. This can take up to 30 minutes..."
  ./build.py

popd

echo "Running tests..."
go test -v .

