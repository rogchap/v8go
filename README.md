# Execute JavaScript from Go

<a href="https://github.com/rogchap/v8go/releases"><img src="https://img.shields.io/github/v/release/rogchap/v8go" alt="Github release"></a>
[![Go Report Card](https://goreportcard.com/badge/rogchap.com/v8go)](https://goreportcard.com/report/rogchap.com/v8go)
[![Go Reference](https://pkg.go.dev/badge/rogchap.com/v8go.svg)](https://pkg.go.dev/rogchap.com/v8go)
[![CI](https://github.com/rogchap/v8go/workflows/CI/badge.svg)](https://github.com/rogchap/v8go/actions?query=workflow%3ACI)
![V8 Build](https://github.com/rogchap/v8go/workflows/V8%20Build/badge.svg)
[![codecov](https://codecov.io/gh/rogchap/v8go/branch/master/graph/badge.svg?token=VHZwzGm3dV)](https://codecov.io/gh/rogchap/v8go)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B22862%2Fgit%40github.com%3Arogchap%2Fv8go.git.svg?type=shield)](https://app.fossa.com/projects/custom%2B22862%2Fgit%40github.com%3Arogchap%2Fv8go.git?ref=badge_shield)
[![#v8go Slack Channel](https://img.shields.io/badge/slack-%23v8go-4A154B?logo=slack)](https://gophers.slack.com/channels/v8go)

<img src="gopher.jpg" width="200px" alt="V8 Gopher based on original artwork from the amazing Renee French" />

## Usage

```go
import v8 "rogchap.com/v8go"
```

### Running a script

```go
ctx := v8.NewContext() // creates a new V8 context with a new Isolate aka VM
ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
ctx.RunScript("const result = add(3, 4)", "main.js") // any functions previously added to the context can be called
val, _ := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
fmt.Printf("addition result: %s", val)
```

### One VM, many contexts

```go
iso := v8.NewIsolate() // creates a new JavaScript VM
ctx1 := v8.NewContext(iso) // new context within the VM
ctx1.RunScript("const multiply = (a, b) => a * b", "math.js")

ctx2 := v8.NewContext(iso) // another context on the same VM
if _, err := ctx2.RunScript("multiply(3, 4)", "main.js"); err != nil {
  // this will error as multiply is not defined in this context
}
```

### JavaScript function with Go callback

```go
iso := v8.NewIsolate() // create a new VM
// a template that represents a JS function
printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
    fmt.Printf("%v", info.Args()) // when the JS function is called this Go callback will execute
    return nil // you can return a value back to the JS caller if required
})
global := v8.NewObjectTemplate(iso) // a template that represents a JS Object
global.Set("print", printfn) // sets the "print" property of the Object to our function
ctx := v8.NewContext(iso, global) // new Context with the global Object set to our object template
ctx.RunScript("print('foo')", "print.js") // will execute the Go callback with a single argunent 'foo'
```

### Update a JavaScript object from Go

```go
ctx := v8.NewContext() // new context with a default VM
obj := ctx.Global() // get the global object from the context
obj.Set("version", "v1.0.0") // set the property "version" on the object
val, _ := ctx.RunScript("version", "version.js") // global object will have the property set within the JS VM
fmt.Printf("version: %s", val)

if obj.Has("version") { // check if a property exists on the object
    obj.Delete("version") // remove the property from the object
}
```

### JavaScript errors

```go
val, err := ctx.RunScript(src, filename)
if err != nil {
  e := err.(*v8.JSError) // JavaScript errors will be returned as the JSError struct
  fmt.Println(e.Message) // the message of the exception thrown
  fmt.Println(e.Location) // the filename, line number and the column where the error occured
  fmt.Println(e.StackTrace) // the full stack trace of the error, if available

  fmt.Printf("javascript error: %v", e) // will format the standard error message
  fmt.Printf("javascript stack trace: %+v", e) // will format the full error stack trace
}
```

### Pre-compile context-independent scripts to speed-up execution times

For scripts that are large or are repeatedly run in different contexts,
it is beneficial to compile the script once and used the cached data from that
compilation to avoid recompiling every time you want to run it.

```go
source := "const multiply = (a, b) => a * b"
iso1 := v8.NewIsolate() // creates a new JavaScript VM
ctx1 := v8.NewContext(iso1) // new context within the VM
script1, _ := iso1.CompileUnboundScript(source, "math.js", v8.CompileOptions{}) // compile script to get cached data
val, _ := script1.Run(ctx1)

cachedData := script1.CreateCodeCache()

iso2 := v8.NewIsolate() // create a new JavaScript VM
ctx2 := v8.NewContext(iso2) // new context within the VM

script2, _ := iso2.CompileUnboundScript(source, "math.js", v8.CompileOptions{CachedData: cachedData}) // compile script in new isolate with cached data
val, _ = script2.Run(ctx2)
```

### Terminate long running scripts

```go
vals := make(chan *v8.Value, 1)
errs := make(chan error, 1)

go func() {
    val, err := ctx.RunScript(script, "forever.js") // exec a long running script
    if err != nil {
        errs <- err
        return
    }
    vals <- val
}()

select {
case val := <- vals:
    // success
case err := <- errs:
    // javascript error
case <- time.After(200 * time.Milliseconds):
    vm := ctx.Isolate() // get the Isolate from the context
    vm.TerminateExecution() // terminate the execution
    err := <- errs // will get a termination error back from the running script
}
```

### CPU Profiler

```go
func createProfile() {
	iso := v8.NewIsolate()
	ctx := v8.NewContext(iso)
	cpuProfiler := v8.NewCPUProfiler(iso)

	cpuProfiler.StartProfiling("my-profile")

	ctx.RunScript(profileScript, "script.js") # this script is defined in cpuprofiler_test.go
	val, _ := ctx.Global().Get("start")
	fn, _ := val.AsFunction()
	fn.Call(ctx.Global())

	cpuProfile := cpuProfiler.StopProfiling("my-profile")

	printTree("", cpuProfile.GetTopDownRoot()) # helper function to print the profile
}

func printTree(nest string, node *v8.CPUProfileNode) {
	fmt.Printf("%s%s %s:%d:%d\n", nest, node.GetFunctionName(), node.GetScriptResourceName(), node.GetLineNumber(), node.GetColumnNumber())
	count := node.GetChildrenCount()
	if count == 0 {
		return
	}
	nest = fmt.Sprintf("%s  ", nest)
	for i := 0; i < count; i++ {
		printTree(nest, node.GetChild(i))
	}
}

// Output
// (root) :0:0
//   (program) :0:0
//   start script.js:23:15
//     foo script.js:15:13
//       delay script.js:12:15
//         loop script.js:1:14
//       bar script.js:13:13
//         delay script.js:12:15
//           loop script.js:1:14
//       baz script.js:14:13
//         delay script.js:12:15
//           loop script.js:1:14
//   (garbage collector) :0:0
```

## Documentation

Go Reference & more examples: https://pkg.go.dev/rogchap.com/v8go

### Support

If you would like to ask questions about this library or want to keep up-to-date with the latest changes and releases,
please join the [**#v8go**](https://gophers.slack.com/channels/v8go) channel on Gophers Slack. [Click here to join the Gophers Slack community!](https://invite.slack.golangbridge.org/)

### Windows

There used to be Windows binary support. For further information see, [this PR](UPDATE).

The v8go library would welcome contributions from anyone able to get an external windows
build of the V8 library linking with v8go, using the version of V8 checked out in the
`deps/v8` git submodule, and documentation of the process involved. This process will likely
involve passing a linker flag when building v8go (e.g. using the `CGO_LDFLAGS` environment
variable.

## V8 dependency

V8 version: **9.0.257.18** (April 2021)

In order to make `v8go` usable as a standard Go package, prebuilt static libraries of V8
are included for Linux and macOS. you *should not* require to build V8 yourself.

Due to security concerns of binary blobs hiding malicious code, the V8 binary is built via CI *ONLY*.

## Project Goals

To provide a high quality, idiomatic, Go binding to the [V8 C++ API](https://v8.github.io/api/head/index.html).

The API should match the original API as closely as possible, but with an API that Gophers (Go enthusiasts) expect. For
example: using multiple return values to return both result and error from a function, rather than throwing an
exception.

This project also aims to keep up-to-date with the latest (stable) release of V8.

## License

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B22862%2Fgit%40github.com%3Arogchap%2Fv8go.git.svg?type=large)](https://app.fossa.com/projects/custom%2B22862%2Fgit%40github.com%3Arogchap%2Fv8go.git?ref=badge_large)

## Development

### Recompile V8 with debug info and debug checks

[Aside from data races, Go should be memory-safe](https://research.swtch.com/gorace) and v8go should preserve this property by adding the necessary checks to return an error or panic on these unsupported code paths. Release builds of v8go don't include debugging information for the V8 library since it significantly adds to the binary size, slows down compilation and shouldn't be needed by users of v8go. However, if a v8go bug causes a crash (e.g. during new feature development) then it can be helpful to build V8 with debugging information to get a C++ backtrace with line numbers. The following steps will not only do that, but also enable V8 debug checking, which can help with catching misuse of the V8 API.

1) Make sure to clone the projects submodules (ie. the V8's `depot_tools` project): `git submodule update --init --recursive`
1) Build the V8 binary for your OS: `deps/build.py --debug`. V8 is a large project, and building the binary can take up to 30 minutes.
1) Build the executable to debug, using `go build` for commands or `go test -c` for tests. You may need to add the `-ldflags=-compressdwarf=false` option to disable debug information compression so this information can be read by the debugger (e.g. lldb that comes with Xcode v12.5.1, the latest Xcode released at the time of writing)
1) Run the executable with a debugger (e.g. `lldb -- ./v8go.test -test.run TestThatIsCrashing`, `run` to start execution then use `bt` to print a bracktrace after it breaks on a crash), since backtraces printed by Go or V8 don't currently include line number information.

### Upgrading the V8 binaries

We have the [upgradev8](https://github.com/rogchap/v8go/.github/workflow/v8upgrade.yml) workflow.
The workflow is triggered every day or manually.

If the current [v8_version](https://github.com/rogchap/v8go/deps/v8_version) is different from the latest stable version, the workflow takes care of fetching the latest stable v8 files and copying them into `deps/include`. The last step of the workflow opens a new PR with the branch name `v8_upgrade/<v8-version>` with all the changes.

The next steps are:

1) The build is not yet triggered automatically. To trigger it manually, go to the [V8
Build](https://github.com/rogchap/v8go/actions?query=workflow%3A%22V8+Build%22) Github Action, Select "Run workflow",
and select your pushed branch eg. `v8_upgrade/<v8-version>`.
1) Once built, this should open 2 PRs against your branch to add the `libv8.a` for Linux and macOS; merge
these PRs into your branch. You are now ready to raise the PR against `master` with the latest version of V8.

### Flushing after C/C++ standard library printing for debugging

When using the C/C++ standard library functions for printing (e.g. `printf`), then the output will be buffered by default.
This can cause some confusion, especially because the test binary (created through `go test`) does not flush the buffer
at exit (at the time of writing). When standard output is the terminal, then it will use line buffering and flush when
a new line is printed, otherwise (e.g. if the output is redirected to a pipe or file) it will be fully buffered and not even
flush at the end of a line. When the test binary is executed through `go test .` (e.g. instead of
separately compiled with `go test -c` and run with `./v8go.test`) Go may redirect standard output internally, resulting in
standard output being fully buffered.

A simple way to avoid this problem is to flush the standard output stream after printing with the `fflush(stdout);` statement.
Not relying on the flushing at exit can also help ensure the output is printed before a crash.

### Local leak checking

Leak checking is automatically done in CI, but it can be useful to do locally to debug leaks.

Leak checking is done using the [Leak Sanitizer](https://clang.llvm.org/docs/LeakSanitizer.html) which
is a part of LLVM. As such, compiling with clang as the C/C++ compiler seems to produce more complete
backtraces (unfortunately still only of the system stack at the time of writing).

For instance, on a Debian-based Linux system, you can use `sudo apt-get install clang-12` to install a
recent version of clang.  Then CC and CXX environment variables are needed to use that compiler. With
that compiler, the tests can be run as follows

```
CC=clang-12 CXX=clang++-12 go test -c --tags leakcheck && ./v8go.test
```

The separate compile and link commands are currently needed to get line numbers in the backtrace.

On macOS, leak checking isn't available with the version of clang that comes with Xcode, so a separate
compiler installation is needed.  For example, with homebrew, `brew install llvm` will install a version
of clang with support for this. The ASAN_OPTIONS environment variable will also be needed to run the code
with leak checking enabled, since it isn't enabled by default on macOS. E.g. with the homebrew
installation of llvm, the tests can be run with

```
CXX=/usr/local/opt/llvm/bin/clang++ CC=/usr/local/opt/llvm/bin/clang go test -c --tags leakcheck -ldflags=-compressdwarf=false
ASAN_OPTIONS=detect_leaks=1 ./v8go.test
```

The `-ldflags=-compressdwarf=false` is currently (with clang 13) needed to get line numbers in the backtrace.

### Formatting

Go has `go fmt`, C has `clang-format`. Any changes to the `v8go.h|cc` should be formated with `clang-format` with the
"Chromium" Coding style. This can be done easily by running the `go generate` command.

`brew install clang-format` to install on macOS.

---

V8 Gopher image based on original artwork from the amazing [Renee French](http://reneefrench.blogspot.com).
