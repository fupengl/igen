package conf

import "testing"

func TestParseLine(t *testing.T) {
	kv := map[string]string{}

	data := [][2]string{
		{"", ""},
		{"1", "1"},
		{"123456", "123456"},
		{`foo = "this is foo"`, `foo = "this is foo"`},
		{`path = "{{goroot | GOROOT}}"`, `path = "/usr/local/go"`},
		{`stub = "{{no_env | NO_EXIST_ENV}}"`, `stub = "no_env"`},
		{`{{aaa}}`, `aaa`},
		{`{{aaa`, `{{aaa`},
		{`aaa}}`, `aaa}}`},
		{`{aaa}`, `{aaa}`},
		{`{{ {{aaa}}`, `{{ aaa`},
		{`{{ {{aaa|bbb}}`, `{{ aaa`},
		{`new_foo = "{{$foo}}"`, `new_foo = "this is foo"`},
		{`the_foo = "{{foo}}"`, `the_foo = "foo"`},
		{`any = "GOROOT at {{|GOROOT}}, Workspace at {{/home/foo/workspace|GOWORKSPACE_NOT_EXISTS}}"`, `any = "GOROOT at /usr/local/go, Workspace at /home/foo/workspace"`},
		{`count = 3`, `count = 3`},
		{`count = 3"`, `count = 3"`},
		{`count = "3`, `count = "3`},
	}

	for _, d := range data {
		got := parseLine(d[0], kv)
		if got != d[1] {
			t.Fatalf("%s after parse, expect %s, but got: %s", d[0], d[1], got)
		}
	}
}
