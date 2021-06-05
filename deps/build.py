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
exclude_unwind_tables=true
"""

def v8deps():
    spec = "solutions = %s" % gclient_sln
    env = os.environ.copy()
    env["PATH"] = tools_path + os.pathsep + env["PATH"]
    subprocess.check_call(cmd(["gclient", "sync", "--spec", spec]),
                        cwd=deps_path,
                        env=env)

def cmd(args):
    return ["cmd", "/c"] + args if is_windows else args

def os_arch():
    u = platform.uname()
    # "x86_64" is called "amd64" on Windows
    return (u[0] + "_" + u[4]).lower().replace("amd64", "x86_64")

def apply_mingw_patches():
    v8_build_path = os.path.join(v8_path, "build")
    apply_patch("0000-add-mingw-main-code-changes", v8_path)
    apply_patch("0001-add-mingw-toolchain", v8_build_path)
    update_last_change()
    zlib_path = os.path.join(v8_path, "third_party", "zlib")
    zlib_src_gn = os.path.join(deps_path, os_arch(), "zlib.gn")
    zlib_dst_gn = os.path.join(zlib_path, "BUILD.gn")
    shutil.copy(zlib_src_gn, zlib_dst_gn)

def apply_patch(patch_name, working_dir):
    patch_path = os.path.join(deps_path, os_arch(), patch_name + ".patch")
    subprocess.check_call(["git", "apply", "-v", patch_path], cwd=working_dir)

def update_last_change():
    import v8.build.util.lastchange as lastchange
    out_path = os.path.join(v8_path, "build", "util", "LASTCHANGE")
    lastchange.main(["lastchange", "-o", out_path])

def main():
    v8deps()
    if is_windows:
        apply_mingw_patches()
    
    gn_path = os.path.join(tools_path, "gn")
    assert(os.path.exists(gn_path))
    ninja_path = os.path.join(tools_path, "ninja" + (".exe" if is_windows else ""))
    assert(os.path.exists(ninja_path))

    build_path = os.path.join(deps_path, ".build", os_arch())
    env = os.environ.copy()

    is_debug = 'true' if args.debug else 'false'
    is_clang = 'true' if args.clang else 'false'
    gnargs = gn_args % (is_debug, is_clang)
    gen_args = gnargs.replace('\n', ' ')
    
    subprocess.check_call(cmd([gn_path, "gen", build_path, "--args=" + gen_args]),
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
