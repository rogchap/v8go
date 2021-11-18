# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Support for calling constructors functions with NewInstance on Function
- Access "this" from function callback
- value.SameValue(otherValue) function to compare values for sameness
- Undefined, Null functions to get these constant values for the isolate
- Support for calling a method on an object.
- Support for calling `IsExecutionTerminating` on isolate to check if execution is still terminating.
- Support for setting and getting internal fields for template object instances
- Support for CPU profiling
- Add V8 build for Apple Silicon
- Add support for throwing an exception directly via the isolate's ThrowException function.
- Support for compiling a context-dependent UnboundScript which can be run in any context of the isolate it was compiled in.
- Support for creating a code cache from an UnboundScript which can be used to create an UnboundScript in other isolates
to run a pre-compiled script in new contexts.

### Changed
- Removed error return value from NewIsolate which never fails
- Removed error return value from NewContext which never fails
- Removed error return value from Context.Isolate() which never fails
- Removed error return value from NewObjectTemplate and NewFunctionTemplate. Panic if given a nil argument.
- Function Call accepts receiver as first argument.
- Removed Windows support until its build issues are addressed.

### Fixed
- Add some missing error propagation
- Fix crash from template finalizer releasing V8 data, let it be disposed with the isolate
- Fix crash by keeping alive the template while its C++ pointer is still being used
- Fix crash from accessing function template callbacks outside of `RunScript`, such as in `JSONStringify`

## [v0.6.0] - 2021-05-11

### Added
- Promise resolver and promise result
- Convert a Value to a Function and invoke it. Thanks to [@robfig](https://github.com/robfig)
- Windows static binary. Thanks to [@cleiner](https://github.com/cleiner)
- Setting/unsetting of V8 feature flags
- Register promise callbacks in Go. Thanks to [@robfig](https://github.com/robfig)
- Get Function from a template for a given context. Thanks to [@robfig](https://github.com/robfig)

### Changed
- Upgrade to V8 9.0.257.18

### Fixed
- Go GC attempting to free C memory (via finalizer) of values after an Isolate is disposed causes a panic

## [v0.5.1] - 2021-02-19

### Fixed
- Memory being held by Values after the associated Context is closed

## [v0.5.0] - 2021-02-08

### Added
- Support for the BigInt value to the big.Int Go type
- Create Object Templates with primitive values, including other Object Templates
- Configure Object Template as the global object of any new Context
- Function Templates with callbacks to Go
- Value to Object type, including Get/Set/Has/Delete methods
- Get Global Object from the Context
- Convert an Object Template to an instance of an Object

### Changed
- NewContext() API has been improved to handle optional global object, as well as optional Isolate
- Package error messages are now prefixed with `v8go` rather than the struct name
- Deprecated `iso.Close()` in favor of `iso.Dispose()` to keep consistancy with the C++ API
- Upgraded V8 to 8.8.278.14
- Licence BSD 3-Clause (same as V8 and Go)

## [v0.4.0] - 2021-01-14

### Added
- Value methods for checking value kind (is string, number, array etc)
- C formatting via `clang-format` to aid future development
- Support of vendoring with `go mod vendor`
- Value methods to convert to primitive data types

### Changed
- Use g++ (default for cgo) for linux builds of the static v8 lib

## [v0.3.0] - 2020-12-18

### Added
- Support for Windows via [MSYS2](https://www.msys2.org/). Thanks to [@neptoess](https://github.com/neptoess)

### Changed
- Upgraded V8 to 8.7.220.31

## [v0.2.0] - 2020-01-25

### Added
- Manually dispose of the isolate when required
- Monitor isolate heap statistics. Thanks to [@mehrdadrad](https://github.com/mehrdadrad)

### Changed
- Upgrade V8 to 8.0.426.15

## [v0.1.0] - 2019-09-22

### Changed
- Upgrade V8 to 7.7.299.9

## [v0.0.1] - 2019-09-2020

### Added
- Create V8 Isolate
- Create Contexts
- Run JavaScript scripts
- Get Values back from JavaScript in Go
- Get detailed JavaScript errors in Go, including stack traces
- Terminate long running scripts from any Goroutine
