// Copyright 2021 the v8go contributors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package v8go_test

import (
	"testing"

	v8 "rogchap.com/v8go"
)

func TestCreateSnapshot(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	snapshotCreatorIso, err := snapshotCreator.GetIsolate()
	fatalIf(t, err)

	snapshotCreatorCtx := v8.NewContext(snapshotCreatorIso)
	defer snapshotCreatorCtx.Close()

	snapshotCreatorCtx.RunScript(`const add = (a, b) => a + b`, "add.js")
	snapshotCreatorCtx.RunScript(`function run() { return add(3, 4); }`, "main.js")
	err = snapshotCreator.SetDefaultContext(snapshotCreatorCtx)
	fatalIf(t, err)

	data, err := snapshotCreator.Create(v8.FunctionCodeHandlingClear)
	fatalIf(t, err)

	iso := v8.NewIsolate(v8.WithStartupData(data))
	defer iso.Dispose()

	ctx := v8.NewContext(iso)
	defer ctx.Close()

	runVal, err := ctx.Global().Get("run")
	fatalIf(t, err)

	fn, err := runVal.AsFunction()
	fatalIf(t, err)

	val, err := fn.Call(v8.Undefined(iso))
	fatalIf(t, err)

	if val.String() != "7" {
		t.Fatal("invalid val")
	}
}

func TestCreateSnapshotAndAddExtraContext(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	snapshotCreatorIso, err := snapshotCreator.GetIsolate()
	fatalIf(t, err)

	snapshotCreatorCtx := v8.NewContext(snapshotCreatorIso)
	defer snapshotCreatorCtx.Close()

	snapshotCreatorCtx.RunScript(`const add = (a, b) => a + b`, "add.js")
	snapshotCreatorCtx.RunScript(`function run() { return add(3, 4); }`, "main.js")
	err = snapshotCreator.SetDefaultContext(snapshotCreatorCtx)
	fatalIf(t, err)

	snapshotCreatorCtx2 := v8.NewContext(snapshotCreatorIso)
	defer snapshotCreatorCtx2.Close()

	snapshotCreatorCtx2.RunScript(`const multiply = (a, b) => a * b`, "add.js")
	snapshotCreatorCtx2.RunScript(`function run() { return multiply(3, 4); }`, "main.js")
	index, err := snapshotCreator.AddContext(snapshotCreatorCtx2)
	fatalIf(t, err)

	snapshotCreatorCtx3 := v8.NewContext(snapshotCreatorIso)
	defer snapshotCreatorCtx3.Close()

	snapshotCreatorCtx3.RunScript(`const div = (a, b) => a / b`, "add.js")
	snapshotCreatorCtx3.RunScript(`function run() { return div(6, 2); }`, "main.js")
	index2, err := snapshotCreator.AddContext(snapshotCreatorCtx3)
	fatalIf(t, err)

	data, err := snapshotCreator.Create(v8.FunctionCodeHandlingClear)
	fatalIf(t, err)

	iso := v8.NewIsolate(v8.WithStartupData(data))
	defer iso.Dispose()

	ctx, err := v8.NewContextFromSnapshot(iso, index)
	fatalIf(t, err)
	defer ctx.Close()

	runVal, err := ctx.Global().Get("run")
	fatalIf(t, err)

	fn, err := runVal.AsFunction()
	fatalIf(t, err)

	val, err := fn.Call(v8.Undefined(iso))
	fatalIf(t, err)

	if val.String() != "12" {
		t.Fatal("invalid val")
	}

	ctx, err = v8.NewContextFromSnapshot(iso, index2)
	fatalIf(t, err)

	defer ctx.Close()

	runVal, err = ctx.Global().Get("run")
	fatalIf(t, err)

	fn, err = runVal.AsFunction()
	fatalIf(t, err)

	val, err = fn.Call(v8.Undefined(iso))
	fatalIf(t, err)

	if val.String() != "3" {
		t.Fatal("invalid val")
	}
}

func TestCreateSnapshotErrorAfterAddingMultipleDefaultContext(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	defer snapshotCreator.Dispose()
	snapshotCreatorIso, err := snapshotCreator.GetIsolate()
	fatalIf(t, err)
	snapshotCreatorCtx := v8.NewContext(snapshotCreatorIso)
	defer snapshotCreatorCtx.Close()

	snapshotCreatorCtx.RunScript(`const add = (a, b) => a + b`, "add.js")
	snapshotCreatorCtx.RunScript(`function run() { return add(3, 4); }`, "main.js")
	err = snapshotCreator.SetDefaultContext(snapshotCreatorCtx)
	fatalIf(t, err)

	err = snapshotCreator.SetDefaultContext(snapshotCreatorCtx)

	if err == nil {
		t.Error("setting another default context should have failed, got <nil>")
	}
}

func TestCreateSnapshotErrorAfterSuccessfullCreate(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	snapshotCreatorIso, err := snapshotCreator.GetIsolate()
	fatalIf(t, err)
	snapshotCreatorCtx := v8.NewContext(snapshotCreatorIso)
	defer snapshotCreatorCtx.Close()

	snapshotCreatorCtx.RunScript(`const add = (a, b) => a + b`, "add.js")
	snapshotCreatorCtx.RunScript(`function run() { return add(3, 4); }`, "main.js")
	err = snapshotCreator.SetDefaultContext(snapshotCreatorCtx)
	fatalIf(t, err)

	_, err = snapshotCreator.Create(v8.FunctionCodeHandlingClear)
	fatalIf(t, err)

	_, err = snapshotCreator.GetIsolate()
	if err == nil {
		t.Error("getting Isolate should have fail")
	}

	_, err = snapshotCreator.AddContext(snapshotCreatorCtx)
	if err == nil {
		t.Error("adding context should have fail")
	}

	_, err = snapshotCreator.Create(v8.FunctionCodeHandlingClear)
	if err == nil {
		t.Error("creating snapshot should have fail")
	}
}

func TestCreateSnapshotErrorIfNoDefaultContextIsAdded(t *testing.T) {
	snapshotCreator := v8.NewSnapshotCreator()
	defer snapshotCreator.Dispose()

	_, err := snapshotCreator.Create(v8.FunctionCodeHandlingClear)

	if err == nil {
		t.Error("creating a snapshop should have fail")
	}
}
