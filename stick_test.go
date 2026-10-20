package stick

import (
	"bytes"
	"testing"
)

func TestExecuteTree(t *testing.T) {
	env := New(nil)
	tree, err := env.Parse("Hello, {{ name }}!")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"World", "Tyler"} {
		var buf bytes.Buffer
		if err := env.ExecuteTree(tree, &buf, map[string]Value{"name": name}); err != nil {
			t.Fatal(err)
		}
		if want := "Hello, " + name + "!"; buf.String() != want {
			t.Errorf("got %q, want %q", buf.String(), want)
		}
	}
}
