#!/usr/bin/env python
import platform
import os
import subprocess

deps_path = os.path.dirname(os.path.realpath(__file__))
v8_path = os.path.join(deps_path, "v8")
tools_path = os.path.join(deps_path, "depot_tools")
is_windows = platform.system().lower() == "windows"

gclient_sln = [
    { "name"        : "v8",
        "url"         : "https://chromium.googlesource.com/v8/v8.git",
        "deps_file"   : "DEPS",
        "managed"     : False,
        "custom_deps" : {
            # These deps are unnecessary for building.
            "v8/testing/gmock"                      : None,
            "v8/test/wasm-js"                       : None,
            "v8/third_party/android_tools"          : None,
            "v8/third_party/catapult"               : None,
            "v8/third_party/colorama/src"           : None,
            "v8/tools/gyp"                          : None,
            "v8/tools/luci-go"                      : None,
        },
        "custom_vars": {
            "build_for_node" : True,
        },
    },
]

def v8deps():
    spec = "solutions = %s" % gclient_sln
    env = os.environ.copy()
    env["PATH"] = tools_path + os.pathsep + env["PATH"]
    subprocess.check_call(cmd(["gclient", "sync", "--spec", spec]),
                        cwd=deps_path,
                        env=env)

def cmd(args):
    return ["cmd", "/c"] + args if is_windows else args


if __name__ == "__main__":
    v8deps()
