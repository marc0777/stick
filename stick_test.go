package stick

import (
	"bytes"
	"errors"
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

// Strict mode follows PHP Twig strict_variables: only names actually read must exist.
func TestStrict(t *testing.T) {
	env := New(nil)
	env.Strict = true
	env.Filters["default"] = func(_ Context, val Value, args ...Value) Value {
		if val == nil && len(args) > 0 {
			return args[0]
		}
		return val
	}
	run := func(tpl string, ctx map[string]Value) (string, error) {
		var buf bytes.Buffer
		err := env.Execute(tpl, &buf, ctx)
		return buf.String(), err
	}
	for tpl, want := range map[string]string{
		"{{ x|default('d') }}": "d",
		"{{ o.a|default('d') }}{{ o.a.b|default('d') }}{{ x.a|default('d') }}":        "ddd",
		"{% if x is defined %}{{ x }}{% else %}n{% endif %}":                          "n",
		"{% if false %}{{ x }}{% endif %}ok":                                          "ok",
		"{% set y = 1 %}{{ y }}{% for i in [2] %}{{ i }}{{ loop.index }}{% endfor %}": "121",
	} {
		if out, err := run(tpl, map[string]Value{"o": map[string]Value{}}); err != nil || out != want {
			t.Errorf("%q: got %q, %v; want %q", tpl, out, err, want)
		}
	}
	for _, tpl := range []string{"{{ x }}", "{% if x %}{% endif %}", "{{ x.y }}", "{% if o[k] is defined %}{% endif %}", "{{ o.a|default(z) }}"} {
		_, err := run(tpl, map[string]Value{"o": map[string]Value{}})
		if undef := (UndefinedVariableError{}); !errors.As(err, &undef) {
			t.Errorf("%q: got %v, want UndefinedVariableError", tpl, err)
		}
	}
	env.Strict = false
	if out, err := run("{{ x }}", nil); err != nil || out != "" {
		t.Errorf("non-strict: got %q, %v", out, err)
	}
}
