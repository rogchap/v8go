// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Ignore leaks within Go standard libraries http/https support code.
// The getaddrinfo detected leaks can be avoided using GODEBUG=netdns=go but
// currently there are more for loading system root certificates on macOS.
//go:build !leakcheck || !darwin
// +build !leakcheck !darwin

package v8go_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	v8 "rogchap.com/v8go"
)

func ExampleFunctionTemplate_fetch() {
	iso := v8.NewIsolate()
	defer iso.Dispose()
	global := v8.NewObjectTemplate(iso)

	fetchfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		args := info.Args()
		url := args[0].String()

		resolver, _ := v8.NewPromiseResolver(info.Context())

		go func() {
			res, _ := http.Get(url)
			body, _ := ioutil.ReadAll(res.Body)
			val := v8.MustNewString(iso, string(body))
			resolver.Resolve(val)
		}()
		return resolver.GetPromise().Value
	})
	global.Set("fetch", fetchfn, v8.ReadOnly)

	ctx := v8.NewContext(iso, global)
	defer ctx.Close()
	val, _ := ctx.RunScript("fetch('https://rogchap.com/v8go')", "")
	prom, _ := val.AsPromise()

	// wait for the promise to resolve
	for prom.State() == v8.Pending {
		continue
	}
	fmt.Printf("%s\n", strings.Split(prom.Result().String(), "\n")[0])
	// Output:
	// <!DOCTYPE html>
}
