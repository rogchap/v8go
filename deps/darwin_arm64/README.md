# macOS arm64

The `arm64` builds are performed locally since GHA doesn't support hosted ARM runners.

To run manually, connect to the Airplane Tailscale network then `ssh` onto the M1 Mini:

```sh
ssh colin@<ip>
```

Kick off a build (would recommend using `screen`):

```sh
cd ~/dev/airplanedev/v8go/deps
git pull
./build.py
```

That takes 30-60 minutes. Once that finishes, run the following on your computer to copy over the binary:

```sh
scp colin@<ip>:/Users/colin/dev/airplanedev/v8go/deps/darwin_arm64/libv8.a ./deps/darwin_arm64/libv8.a
```

Finally, open a PR and merge!
