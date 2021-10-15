// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestBuiltinSymbol(t *testing.T) {
	t.Parallel()

	iso := v8.NewIsolate()
	defer iso.Dispose()

	tsts := []struct {
		Func            func(*v8.Isolate) *v8.Symbol
		WantDescription string
	}{
		{v8.SymbolAsyncIterator, "Symbol.asyncIterator"},
		{v8.SymbolHasInstance, "Symbol.hasInstance"},
		{v8.SymbolIsConcatSpreadable, "Symbol.isConcatSpreadable"},
		{v8.SymbolIterator, "Symbol.iterator"},
		{v8.SymbolMatch, "Symbol.match"},
		{v8.SymbolReplace, "Symbol.replace"},
		{v8.SymbolSearch, "Symbol.search"},
		{v8.SymbolSplit, "Symbol.split"},
		{v8.SymbolToPrimitive, "Symbol.toPrimitive"},
		{v8.SymbolToStringTag, "Symbol.toStringTag"},
		{v8.SymbolUnscopables, "Symbol.unscopables"},
	}
	for _, tst := range tsts {
		t.Run(tst.WantDescription, func(t *testing.T) {
			iter := tst.Func(iso)
			if iter.Description() != tst.WantDescription {
				t.Errorf("Description: got %q, want %q", iter.Description(), tst.WantDescription)
			}
		})
	}
}
