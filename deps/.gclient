solutions = [
  {
    "name"       : "v8",
    "url"        : "https://chromium.googlesource.com/v8/v8.git",
    "deps_file"  : "DEPS",
    "managed"    : False,
    "custom_deps": {
      # These deps are unnecessary for building.
      #
      # https://chromium.googlesource.com/v8/v8/+/refs/heads/master/DEPS
      # https://github.com/nodejs/node/blob/e83c8aef4b787b531f23ad7abb498fe81db5ee83/tools/v8/fetch_deps.py
      # https://github.com/rogchap/v8go/blob/master/deps/build.py
      "v8/test/benchmarks/data"               : None,
      "v8/testing/gmock"                      : None,
      "v8/test/wasm-js"                       : None,
      "v8/test/mozilla/data"                  : None,
      "v8/test/test262/data"                  : None,
      "v8/test/test262/harness"               : None,
      "v8/third_party/android_ndk"            : None,
      "v8/third_party/android_sdk"            : None,
      "v8/third_party/catapult"               : None,
      "v8/third_party/colorama/src"           : None,
      "v8/third_party/fuchsia-sdk"            : None,
      "v8/third_party/instrumented_libraries" : None,
      "v8/tools/gyp"                          : None,
      "v8/tools/luci-go"                      : None,
      "v8/tools/swarming_client"              : None,
      "v8/third_party/qemu-linux-x64"         : None,
    },
  },
]