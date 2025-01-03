#!/bin/sh

set -e

current_v8_version="$(head -n1 VERSION | cut -d'-' -f1)"

latest_chromium_version="$(wget -qO- "https://chromiumdash.appspot.com/fetch_releases?channel=Stable&platform=Linux&num=1" | jq -r '.[0].version')"
latest_v8_version="$(wget -qO- "https://chromiumdash.appspot.com/fetch_version?version=$latest_chromium_version" | jq -r '.v8_version')"

echo "Current version: $current_v8_version"
echo "Latest version: $latest_v8_version"

if [ -n "$latest_v8_version" ] && [ x"$latest_v8_version" != x"$current_v8_version" ] ; then
    echo "Upgrading V8 $current_v8_version --> $latest_v8_version"
    echo -n "$latest_v8_version" > VERSION
    echo "NEW_V8_VERSION=$latest_v8_version" >> $GITHUB_ENV
    echo "NEW_V8_VERSION=$latest_v8_version" >> $GITHUB_OUTPUT
else
    echo "V8 is up to date on version $current_v8_version"
    echo "NEW_V8_VERSION=0" >> $GITHUB_ENV
    echo "NEW_V8_VERSION=0" >> $GITHUB_OUTPUT
fi
