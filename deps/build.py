#!/usr/bin/env python
import platform
import os
import subprocess
import shutil
import argparse

valid_archs = ['arm64', 'x86_64']
# "x86_64" is called "amd64" on Windows
current_arch = platform.uname()[4].lower().replace("amd64", "x86_64")
default_arch = current_arch if current_arch in valid_archs else None

parser = argparse.ArgumentParser()
parser.add_argument('--debug', dest='debug', action='store_true')
parser.add_argument('--no-clang', dest='clang', action='store_false')
parser.add_argument('--arch',
    dest='arch',
    action='store',
    choices=valid_archs,
    default=default_arch,
    required=default_arch is None)
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
target_cpu="%s"
v8_target_cpu="%s"
clang_use_chrome_plugins=false
use_custom_libcxx=false
use_sysroot=false
symbol_level=%s
strip_debug_info=%s
is_component_build=false
v8_monolithic=true
v8_use_external_startup_data=false
treat_warnings_as_errors=false
v8_embedder_string="-v8go"
v8_enable_gdbjit=false
v8_enable_i18n_support=true
icu_use_data_file=false
v8_enable_test_features=false
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
    return u[0].lower() + "_" + args.arch

def v8_arch():
    if args.arch == "x86_64":
        return "x64"
    return args.arch

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
    out_path = os.path.join(v8_path, "build", "util", "LASTCHANGE")
    subprocess.check_call(["python", "build/util/lastchange.py", "-o", out_path], cwd=v8_path)

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
    # symbol_level = 1 includes line number information
    # symbol_level = 2 can be used for additional debug information, but it can increase the
    #   compiled library by an order of magnitude and further slow down compilation
    symbol_level = 1 if args.debug else 0
    strip_debug_info = 'false' if args.debug else 'true'

    arch = v8_arch()
    gnargs = gn_args % (is_debug, is_clang, arch, arch, symbol_level, strip_debug_info)
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
