package v8go

import "testing"

func TestVersion(t *testing.T) {
	print(Version())
	t.Fatal("error")
}
