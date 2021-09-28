#!/usr/bin/python3

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
  _ "rogchap.com/v8go/deps/linux_x86_64"
  _ "rogchap.com/v8go/deps/windows_x86_64"
  _ "rogchap.com/v8go/deps/include"
  %s
)
"""

deps_path = os.path.dirname(os.path.realpath(__file__))
v8go_path = os.path.abspath(os.path.join(deps_path, os.pardir))
v8_include_path = os.path.join(deps_path, "v8", "include")
deps_include_path = os.path.join(deps_path, "include")

def get_directories_names(path):
  flist = []
  for p in pathlib.Path(path).iterdir():
    if p.is_dir():
        flist.append(p.name)
  return flist

def delete_files_except_vendor_files(path):
  for rootDir, subdirs, filenames in os.walk(path):
    # Find the files that matches the given patterm
    for file in filenames:
      if not fnmatch.fnmatch(file, '*.go'):
        try:
            os.remove(os.path.join(rootDir, file))
        except OSError:
            print("Error while deleting file")

def copy_include_directory(src_path, dest_path):
  shutil.copytree(src_path, dest_path, dirs_exist_ok=True)

def update_cgo_file(src_path):
  directories = get_directories_names(src_path)
  package_names = []

  for directory_name in directories:
    package_name = "_ \"rogchap.com/v8go/deps/include/" + directory_name + "\"\n"
    package_names.append(package_name)

  with open(os.path.join(v8go_path, 'cgo.go'), 'w') as temp_file:
      temp_file.write(cgo_file_template % ('  '.join(package_names)))

def create_new_directories_and_vendor_files(directories, src_path):
  for directory in directories:
    new_directory_path = os.path.join(src_path, directory)
    os.mkdir(new_directory_path)

    with open(os.path.join(new_directory_path, 'vendor.go'), 'w') as temp_file:
      temp_file.write(vendor_file_template % (directory, directory))

  update_cgo_file(src_path)


v8_include_directories_names = get_directories_names(v8_include_path)
deps_include_directories_names = get_directories_names(deps_include_path)

if not v8_include_directories_names == deps_include_directories_names:
  new_directories = set(v8_include_directories_names).difference(set(deps_include_directories_names))
  create_new_directories_and_vendor_files(new_directories, deps_include_path)

delete_files_except_vendor_files(deps_include_path)
copy_include_directory(v8_include_path, deps_include_path)
