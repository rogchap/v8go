#!/bin/sh

set -e

target_cpu="$1"

dir="$(cd "$(dirname "$0")" && pwd)"
v8_dir="${dir}/v8"

if [ ! -d "$v8_dir" ]; then
  echo "v8 not found at $v8_dir"
  exit 1
fi

depot_tools_dir="${v8_dir}/third_party/depot_tools"

if [ ! -d "$depot_tools_dir" ]; then
  depot_tools_dir="${dir}/depot_tools"
fi

PATH="${depot_tools_dir}:$PATH"
export PATH

os=""
case "$(uname -s)" in
  Linux)
    os="linux"
    ;;
  Darwin)
    os="darwin"
    ;;
  *)
    echo "Unknown OS type"
    exit 1
esac

cores="2"

if [ "$os" = "linux" ]; then
  cores="$(grep -c processor /proc/cpuinfo)"
elif [ "$os" = "darwin" ]; then
  cores="$(sysctl -n hw.logicalcpu)"
fi

echo "Building V8 for $os $target_cpu"

cc_wrapper=""
if command -v ccache >/dev/null 2>&1 ; then
  cc_wrapper="ccache"
fi

gn_args="$(grep -v "^#" "${dir}/args/${os}.gn" | grep -v "^$")
cc_wrapper=\"$cc_wrapper\"
target_cpu=\"$target_cpu\"
v8_target_cpu=\"$target_cpu\""

cd "${dir}/v8"

gn gen "./out/release" --args="$gn_args"

echo "==================== Build args start ===================="
gn args "./out/release" --list | tee "${dir}/gn-args_${os}.txt"
echo "==================== Build args end ===================="

(
  set -x
  ninja -C "./out/release" -j "$cores" v8_monolith
)

ls -lh ./out/release/obj/libv8_*.a

cd -
