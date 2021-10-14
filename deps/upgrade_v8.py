#!/usr/bin/env python

import urllib.request
import sys
import json
import subprocess
import os
import fnmatch
import pathlib
import shutil

vendor_file_template = """// Package %s is required to provide support for vendoring modules
// DO NOT REMOVE
package %s
"""

cgo_file_template = """// Copyright 2019 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

//go:generate clang-format -i --verbose -style=Chromium v8go.h v8go.cc

// #cgo CXXFLAGS: -fno-rtti -fpic -std=c++14 -DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH -I${SRCDIR}/deps/include
// #cgo LDFLAGS: -pthread -lv8
// #cgo darwin LDFLAGS: -L${SRCDIR}/deps/darwin_x86_64
// #cgo linux LDFLAGS: -L${SRCDIR}/deps/linux_x86_64
// #cgo windows LDFLAGS: -L${SRCDIR}/deps/windows_x86_64 -static -ldbghelp -lssp -lwinmm -lz
import "C"

// These imports forces `go mod vendor` to pull in all the folders that
// contain V8 libraries and headers which otherwise would be ignored.
// DO NOT REMOVE
import (
  _ "rogchap.com/v8go/deps/darwin_x86_64"
  _ "rogchap.com/v8go/deps/include"
  %s
  _ "rogchap.com/v8go/deps/linux_x86_64"
  _ "rogchap.com/v8go/deps/windows_x86_64"
)
"""

CHROME_VERSIONS_URL = "https://omahaproxy.appspot.com/all.json?os=linux&channel=stable"
V8_VERSION_FILE = "v8_version"

deps_path = os.path.dirname(os.path.realpath(__file__))
v8go_path = os.path.abspath(os.path.join(deps_path, os.pardir))
env = os.environ.copy()
v8_path = os.path.join(deps_path, "v8", "include")
v8_include_path = os.path.join(deps_path, "v8", "include")
deps_include_path = os.path.join(deps_path, "include")

def get_directories_names(path):
  flist = []
  for p in pathlib.Path(path).iterdir():
    if p.is_dir():
        flist.append(p.name)
  return sorted(flist)

def delete_files_except_vendor_files(path):
  for rootDir, subdirs, filenames in os.walk(path):
    # Find the files that matches the given patterm
    for file in filenames:
      if not fnmatch.fnmatch(file, '*.go'):
        try:
            os.remove(os.path.join(rootDir, file))
        except OSError:
            print("Error while deleting file")

def package_name(package, index, total):

  name = f'_ "rogchap.com/v8go/deps/include/{package}"'
  if (index + 1 == total):
    return name
  else:
    return name + '\n'

def update_cgo_file(src_path):
  directories = get_directories_names(src_path)
  package_names = []
  total_directories = len(directories)

  for index, directory_name in enumerate(directories):
    package_names.append(package_name(directory_name, index, total_directories))

  with open(os.path.join(v8go_path, 'cgo.go'), 'w') as temp_file:
      temp_file.write(cgo_file_template % ('  '.join(package_names)))

def create_new_vendor_files(src_path):
  directories = get_directories_names(src_path)

  for directory in directories:
    directory_path = os.path.join(src_path, directory)

    vendor_go_file_path = os.path.join(directory_path, 'vendor.go')

    if os.path.isfile(vendor_go_file_path):
      continue

    with open(os.path.join(directory_path, 'vendor.go'), 'w') as temp_file:
      temp_file.write(vendor_file_template % (directory, directory))

  update_cgo_file(src_path)

def update_v8_file(src_path, version):
  with open(os.path.join(src_path, V8_VERSION_FILE), "w") as v8_file:
    v8_file.write(version)

def read_v8_file(src_path):
  v8_version_file = open(os.path.join(src_path, V8_VERSION_FILE), "r")
  return v8_version_file.read().strip()

def get_latest_v8_info():
  with urllib.request.urlopen(CHROME_VERSIONS_URL) as response:
   json_response = response.read()

  return json.loads(json_response)

# Current version
current_v8_version_installed = read_v8_file(deps_path)

# Get latest version
latest_v8_info = get_latest_v8_info()

latest_stable_v8_version = latest_v8_info[0]["versions"][0]["v8_version"]

if not current_v8_version_installed == latest_stable_v8_version:
  # fetch latest v8
  subprocess.run(["fetch", "v8"],
                        cwd=v8_path,
                        env=env)
  subprocess.check_call(["git", "fetch"],
                        cwd=v8_path,
                        env=env)
  # checkout latest stable commit
  subprocess.check_call(["git", "checkout", latest_stable_v8_version],
                        cwd=v8_path,
                        env=env)

  delete_files_except_vendor_files(deps_include_path)
  shutil.copytree(v8_include_path, deps_include_path, dirs_exist_ok=True)
  create_new_vendor_files(deps_include_path)
  update_v8_file(deps_path, latest_stable_v8_version)
