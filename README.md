# Execute JavaScript from Go

[![Go Report Card](https://goreportcard.com/badge/rogchap.com/v8go)](https://goreportcard.com/report/rogchap.com/v8go) 
[![GoDoc](https://godoc.org/rogchap.com/v8go?status.svg)](https://godoc.org/rogchap.com/v8go) 
[![Build Status](https://travis-ci.org/rogchap/v8go.svg?branch=master)](https://travis-ci.org/rogchap/v8go)

![V8 Gopher based on original artwork from the amazing Renee French](gopher.jpg)

## Usage

```go
import "github.com/twitchel/v8go"
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

GoDoc: https://godoc.org/rogchap.com/v8go

## V8 dependency

V8 version: 8.0.426.15

In order to make `v8go` usable as a standard Go package, prebuilt static libraries of V8
are included for Linux and OSX ie. you *should not* require to build V8 yourself.

V8 requires 64-bit, therefore will not work on 32-bit systems. 

---

V8 Gopher image based on original artwork from the amazing [Renee French](http://reneefrench.blogspot.com).
