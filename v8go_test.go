package v8go_test

import (
	"regexp"
	"testing"

	"rogchap.com/v8go"
)

var iso = v8go.NewIsolate()

func TestVersion(t *testing.T) {
	t.Parallel()
	rgx := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+-v8go$`)
	v := v8go.Version()
	if !rgx.MatchString(v) {
		t.Errorf("version string is in the incorrect format: %s", v)
	}
}

func TestRunScriptStringer(t *testing.T) {
	t.Parallel()
	ctx := v8go.NewContext(iso)
	var tests = [...]struct {
		name   string
		source string
		out    string
	}{
		{"Addition", "2 + 2", "4"},
		{"Multiplication", "13 * 2", "26"},
		{"String", `"string"`, "string"},
		{"Object", `let obj = {}; obj`, "[object Object]"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, _ := ctx.RunScript(tt.source)
			str := result.String()
			if str != tt.out {
				t.Errorf("unespected result: expected %q, got %q", tt.out, str)
			}
		})
	}

}
