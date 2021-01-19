# Execute JavaScript from Go

[![Go Report Card](https://goreportcard.com/badge/rogchap.com/v8go)](https://goreportcard.com/report/rogchap.com/v8go) 
[![Go Reference](https://pkg.go.dev/badge/rogchap.com/v8go.svg)](https://pkg.go.dev/rogchap.com/v8go)
[![CI](https://github.com/rogchap/v8go/workflows/CI/badge.svg)](https://github.com/rogchap/v8go/actions?query=workflow%3ACI)
[![#v8go Slack Channel](https://img.shields.io/badge/slack-%23v8go-4A154B?logo=slack)](https://gophers.slack.com/channels/v8go)

<img src="gopher.jpg" width="200px" alt="V8 Gopher based on original artwork from the amazing Renee French" />

## Usage

```go
import "rogchap.com/v8go"
```

### Running a script

```go
ctx, _ := v8go.NewContext(nil) // creates a new V8 context with a new Isolate aka VM
ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
ctx.RunScript("const result = add(3, 4)", "main.js") // any functions previously added to the context can be called
val, _ := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
fmt.Printf("addition result: %s", val)
```

### One VM, many contexts

```go
vm, _ := v8go.NewIsolate() // creates a new JavaScript VM
ctx1, _ := v8go.NewContext(vm) // new context within the VM
ctx1.RunScript("const multiply = (a, b) => a * b", "math.js")

ctx2, _ := v8go.NewContext(vm) // another context on the same VM
if _, err := ctx2.RunScript("multiply(3, 4)", "main.js"); err != nil {
  // this will error as multiply is not defined in this context
}
```

### Javascript errors

```go
val, err := ctx.RunScript(src, filename)
if err != nil {
  err = err.(v8go.JSError) // JavaScript errors will be returned as the JSError struct
  fmt.Println(err.Message) // the message of the exception thrown
  fmt.Println(err.Location) // the filename, line number and the column where the error occured
  fmt.Println(err.StackTrace) // the full stack trace of the error, if available

  fmt.Printf("javascript error: %v", err) // will format the standard error message
  fmt.Printf("javascript stack trace: %+v", err) // will format the full error stack trace
}
```

### Terminate long running scripts

```go

vals := make(chan *v8go.Value, 1)
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
    // sucess
case err := <- errs:
    // javascript error
case <- time.After(200 * time.Milliseconds):
    vm, _ := ctx.Isolate() // get the Isolate from the context
    vm.TerminateExecution() // terminate the execution 
    err := <- errs // will get a termination error back from the running script
}
```

## Documentation

Go Reference: https://pkg.go.dev/rogchap.com/v8go

### Support

If you would like to ask questions about this library or want to keep up-to-date with the latest changes and releases,
please join the [**#v8go**](https://gophers.slack.com/channels/v8go) channel on Gophers Slack. [Click here to join the Gophers Slack community!](https://invite.slack.golangbridge.org/)

## V8 dependency

V8 version: 8.7.220.31

In order to make `v8go` usable as a standard Go package, prebuilt static libraries of V8
are included for Linux and macOS ie. you *should not* require to build V8 yourself.

V8 requires 64-bit, therefore will not work on 32-bit systems. 

## Windows
While no prebuilt static V8 library is included for Windows, MSYS2 provides a package containing
a dynamically linked V8 library that works.

To set this up:
1. Install MSYS2 (https://www.msys2.org/)
2. Add the Mingw-w64 bin to your PATH environment variable (`C:\msys64\mingw64\bin` by default)
3. Open MSYS2 MSYS and execute `pacman -S mingw-w64-x86_64-toolchain mingw-w64-x86_64-v8`
4. This will allow building projects that depend on `v8go`, but, in order to actually run them,
   you will need to copy the `snapshot_blob.bin` file from the Mingw-w64 bin folder to your program's
   working directory (which is typically wherever `main.go` is)

## Development

### Formatting

Go has `go fmt`, C has `clang-format`. Any changes to the `v8go.h|cc` should be formated with `clang-format` with the
"Chromium" Coding style. This can be done easily by running the `go generate` command.

`brew install clang-format` to install on macOS.

---

V8 Gopher image based on original artwork from the amazing [Renee French](http://reneefrench.blogspot.com).
