# v8go provides an API to V8 JavaScript Engine

V8 version: 7.6.303.31

## Usage

```go
import "rogchap.com/v8go"
```

## V8 Dependancy
In order to make `v8go` usable as a standard Go package, prebuilt static libraries of V8
are included for Linux and OSX ie. you should not require to build V8 yourself.

V8 requires 64-bit, therfore will not work on 32-bit systems. 
