package v8go_test

import (
	"regexp"
	"testing"

	"rogchap.com/v8go"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	rgx := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+-v8go$`)
	v := v8go.Version()
	if !rgx.MatchString(v) {
		t.Errorf("version string is in the incorrect format: %s", v)
	}
}

func TestRunScript(t *testing.T) {
	t.Parallel()
	iso := v8go.NewIsolate()
	ctx := v8go.NewContext(iso)
	ctx.RunScript()
	t.Error("See output")
}
