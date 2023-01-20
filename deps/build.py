#!/usr/bin/env python
import platform
import os
import subprocess
import shutil
import argparse
from get_v8deps import v8deps
from compile_v8 import v8compile

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

def main():
    v8deps()
    v8compile(args.debug, args.clang, args.arch)


if __name__ == "__main__":
    main()
