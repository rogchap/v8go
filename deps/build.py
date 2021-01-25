#!/usr/bin/env python
import platform
import os
import subprocess
import shutil
import argparse

parser = argparse.ArgumentParser()
parser.add_argument('--debug', dest='debug', action='store_true')
parser.add_argument('--no-clang', dest='clang', action='store_false')
parser.set_defaults(debug=False, clang=True)
args = parser.parse_args()

deps_path = os.path.dirname(os.path.realpath(__file__))
v8_path = os.path.join(deps_path, "v8")
tools_path = os.path.join(deps_path, "depot_tools")

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

gn_args = """
is_debug=%s
is_clang=%s
clang_use_chrome_plugins=false
use_custom_libcxx=false
use_sysroot=false
symbol_level=0
strip_debug_info=true
is_component_build=false
v8_monolithic=true
v8_use_external_startup_data=false
treat_warnings_as_errors=false
v8_embedder_string="-v8go"
v8_enable_gdbjit=false
v8_enable_i18n_support=false
v8_enable_test_features=false
v8_untrusted_code_mitigations=false
"""

def v8deps():
    spec = "solutions = %s" % gclient_sln
    env = os.environ.copy()
    env["PATH"] = tools_path + os.pathsep + env["PATH"]
    subprocess.check_call(["gclient", "sync", "--spec", spec],
                        cwd=deps_path,
                        env=env)

def os_arch():
    u = platform.uname()
    return (u[0] + "_" + u[4]).lower()

def main():
    v8deps()
    gn_path = os.path.join(tools_path, "gn")
    assert(os.path.exists(gn_path))
    ninja_path = os.path.join(tools_path, "ninja")
    assert(os.path.exists(ninja_path))

    build_path = os.path.join(deps_path, ".build", os_arch())
    env = os.environ.copy()

    is_debug = 'true' if args.debug else 'false'
    is_clang = 'true' if args.clang else 'false'
    gnargs = gn_args % (is_debug, is_clang)
    gen_args = gnargs.replace('\n', ' ')
    
    subprocess.check_call([gn_path, "gen", build_path, "--args=" + gen_args], 
                        cwd=v8_path,
                        env=env)
    subprocess.check_call([ninja_path, "-v", "-C", build_path, "v8_monolith"],
                        cwd=v8_path,
                        env=env)

    lib_fn = os.path.join(build_path, "obj/libv8_monolith.a")
    dest_path = os.path.join(deps_path, os_arch())
    if not os.path.exists(dest_path):
        os.makedirs(dest_path)
    dest_fn = os.path.join(dest_path, 'libv8.a')
    shutil.copy(lib_fn, dest_fn)


if __name__ == "__main__":
    main()
