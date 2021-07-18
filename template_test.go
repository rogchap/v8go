package v8go

import (
	"testing"
)

type templateTest struct {
	*ObjectTemplate
}

func TestUnderlayingTemplater(t *testing.T) {
	iso, _ := NewIsolate()
	ob, _ := NewObjectTemplate(iso)
	_, ok := underlayingTmpl(&templateTest{ObjectTemplate: ob})
	if !ok {
		t.Error("no interface found by `getUnderlayingTemplate`")
	}
}

func TestUnderlayingTemplater2(t *testing.T) {
	iso, _ := NewIsolate()
	ob1, _ := NewObjectTemplate(iso)
	ob2, _ := NewObjectTemplate(iso)

	ob2.Set("hello", &templateTest{ObjectTemplate: ob1})
}
