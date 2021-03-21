
The MINGW patches (`0000`, `0001`, `zlib.gn`) are originally from the [MSYS2 V8 package](https://packages.msys2.org/package/mingw-w64-x86_64-v8?repo=mingw64)
([sources](https://github.com/msys2/MINGW-packages/tree/master/mingw-w64-v8), [LICENSE](https://github.com/msys2/MINGW-packages/blob/master/LICENSE)).

As v8go was already using V8 8.9.255.20 while the MSYS2 package was still at
8.8.278.14-1, they have been slightly updated to work with the newer version.

Patch `0002` has been added to be able to build with `exclude_unwind_tables=true`
(see GN args in `../build.py`). Otherwise the resulting binary exceeds GitHub's
file size limit of 100MB.

To create a new version of the patches from a modified working tree:

    cd ${v8go_project}/deps/v8
    git diff --relative >../windows_x86_64/0000-add-mingw-main-code-changes.patch
    cd build
    git diff --relative >../../windows_x86_64/0001-add-mingw-toolchain.patch
