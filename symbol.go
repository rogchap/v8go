// Copyright 2021 Roger Chapman and the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go

import (
	"fmt"
	"unsafe"

	// #include <stdlib.h>
	// #include "v8go.h"
	"C"
)

// A Symbol represents a JavaScript symbol (ECMA-262 edition 6).
type Symbol struct {
	*Value
}

func SymbolAsyncIterator(iso *Isolate) *Symbol { return symbolByIndex(iso, C.SYMBOL_ASYNC_ITERATOR) }
func SymbolHasInstance(iso *Isolate) *Symbol   { return symbolByIndex(iso, C.SYMBOL_HAS_INSTANCE) }
func SymbolIsConcatSpreadable(iso *Isolate) *Symbol {
	return symbolByIndex(iso, C.SYMBOL_IS_CONCAT_SPREADABLE)
}
func SymbolIterator(iso *Isolate) *Symbol    { return symbolByIndex(iso, C.SYMBOL_ITERATOR) }
func SymbolMatch(iso *Isolate) *Symbol       { return symbolByIndex(iso, C.SYMBOL_MATCH) }
func SymbolReplace(iso *Isolate) *Symbol     { return symbolByIndex(iso, C.SYMBOL_REPLACE) }
func SymbolSearch(iso *Isolate) *Symbol      { return symbolByIndex(iso, C.SYMBOL_SEARCH) }
func SymbolSplit(iso *Isolate) *Symbol       { return symbolByIndex(iso, C.SYMBOL_SPLIT) }
func SymbolToPrimitive(iso *Isolate) *Symbol { return symbolByIndex(iso, C.SYMBOL_TO_PRIMITIVE) }
func SymbolToStringTag(iso *Isolate) *Symbol { return symbolByIndex(iso, C.SYMBOL_TO_STRING_TAG) }
func SymbolUnscopables(iso *Isolate) *Symbol { return symbolByIndex(iso, C.SYMBOL_UNSCOPABLES) }

// symbolByIndex is a Go-to-C helper for obtaining builtin symbols.
func symbolByIndex(iso *Isolate, idx C.SymbolIndex) *Symbol {
	val := C.BuiltinSymbol(iso.ptr, idx)
	if val == nil {
		panic(fmt.Errorf("unknown symbol index: %d", idx))
	}
	return &Symbol{&Value{val, nil}}
}

// Description returns the string representation of the symbol,
// e.g. "Symbol.asyncIterator".
func (sym *Symbol) Description() string {
	s := C.SymbolDescription(sym.Value.ptr)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// String returns Description().
func (sym *Symbol) String() string {
	return sym.Description()
}

// value implements Valuer.
func (sym *Symbol) value() *Value {
	return sym.Value
}
