package v8go_test

import "testing"

func fatalIf(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func recoverPanic(f func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	f()
	return nil
}
