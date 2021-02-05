package v8go_test

import (
	"fmt"
	"testing"

	"rogchap.com/v8go"
)

func TestObjectSet(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	val, _ := ctx.RunScript("const foo = {}; foo", "")
	obj, _ := val.AsObject()
	obj.Set("bar", "baz")
	baz, _ := ctx.RunScript("foo.bar", "")
	if baz.String() != "baz" {
		t.Errorf("unexpected value: %q", baz)
	}
	if err := obj.Set("", nil); err == nil {
		t.Error("expected error but got <nil>")
	}
	if err := obj.Set("a", 0); err == nil {
		t.Error("expected error but got <nil>")
	}
	obj.SetIdx(10, "ten")
	if ten, _ := ctx.RunScript("foo[10]", ""); ten.String() != "ten" {
		t.Errorf("unexpected value: %q", ten)
	}
}

func TestObjectGet(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	val, _ := ctx.RunScript("const foo = { bar: 'baz'}; foo", "")
	obj, _ := val.AsObject()
	if bar, _ := obj.Get("bar"); bar.String() != "baz" {
		t.Errorf("unexpected value: %q", bar)
	}
	if baz, _ := obj.Get("baz"); !baz.IsUndefined() {
		t.Errorf("unexpected value: %q", baz)
	}
	ctx.RunScript("foo[5] = 5", "")
	if five, _ := obj.GetIdx(5); five.Integer() != 5 {
		t.Errorf("unexpected value: %q", five)
	}
	if u, _ := obj.GetIdx(55); !u.IsUndefined() {
		t.Errorf("unexpected value: %q", u)
	}
}

func TestObjectHas(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	val, _ := ctx.RunScript("const foo = {a: 1, '2': 2}; foo", "")
	obj, _ := val.AsObject()
	if !obj.Has("a") {
		t.Error("expected true, got false")
	}
	if obj.Has("c") {
		t.Error("expected false, got true")
	}
	if !obj.HasIdx(2) {
		t.Error("expected true, got false")
	}
	if obj.HasIdx(1) {
		t.Error("expected false, got true")
	}
}

func TestObjectDelete(t *testing.T) {
	t.Parallel()

	ctx, _ := v8go.NewContext()
	val, _ := ctx.RunScript("const foo = { bar: 'baz', '2': 2}; foo", "")
	obj, _ := val.AsObject()
	if !obj.Has("bar") {
		t.Error("expected property to exist")
	}
	if !obj.Delete("bar") {
		t.Error("expected delete to return true, got false")
	}
	if obj.Has("bar") {
		t.Error("expected property to be deleted")
	}
	if !obj.DeleteIdx(2) {
		t.Error("expected delete to return true, got false")
	}

}

func ExampleObject_global() {
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)
	global := ctx.Global()

	console, _ := v8go.NewObjectTemplate(iso)
	logfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Println(info.Args()[0])
		return nil
	})
	console.Set("log", logfn)
	consoleObj, _ := console.NewInstance(ctx)

	global.Set("console", consoleObj)
	ctx.RunScript("console.log('foo')", "")
	// Output:
	// foo
}
