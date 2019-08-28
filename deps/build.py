#!/usr/bin/env python
import platform
import os
import subprocess

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
            "v8/test/benchmarks/data"               : None,
            "v8/testing/gmock"                      : None,
            "v8/test/mozilla/data"                  : None,
            "v8/test/test262/data"                  : None,
            "v8/test/test262/harness"               : None,
            "v8/test/wasm-js"                       : None,
            "v8/third_party/android_tools"          : None,
            "v8/third_party/catapult"               : None,
            "v8/third_party/colorama/src"           : None,
            "v8/third_party/instrumented_libraries" : None,
            "v8/tools/gyp"                          : None,
            "v8/tools/luci-go"                      : None,
            "v8/tools/swarming_client"              : None,
        },
        "custom_vars": {
            "build_for_node" : True,
        },
    },
]

gn_args = """
    is_component_build=false
    is_debug=false
    symbol_level=0
    strip_debug_info=true
    treat_warnings_as_errors=false
    use_custom_libcxx=false
    use_sysroot=false
    v8_deprecation_warnings=false
    v8_embedder_string="-v8worker2"
    v8_enable_gdbjit=false
    v8_enable_i18n_support=false
    v8_enable_test_features=false
    v8_extra_library_files=[]
    v8_imminent_deprecation_warnings=false
    v8_monolithic=true
    v8_static_library=true
    v8_target_cpu="x64"
    v8_untrusted_code_mitigations=false
    v8_use_external_startup_data=false
    v8_use_snapshot=true
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
    return (u[0] + "-" + u[4]).lower()

def main():
    v8deps()
    gn_path = os.path.join(tools_path, "gn")
    assert(os.path.exists(gn_path))
    ninja_path = os.path.join(tools_path, "ninja")
    assert(os.path.exists(ninja_path))

    build_path = os.path.join(deps_path, "build")

    env = os.environ.copy()
    args = gn_args.replace('\n', ' ')
    subprocess.check_call([gn_path, "gen", build_path, "--args=" + args], 
                        cwd=v8_path,
                        env=env)
    subprocess.check_call([ninja_path, "-v", "-C", build_path, "v8_monolith"],
                        cwd=v8_path,
                        env=env)

if __name__ == "__main__":
    main()
