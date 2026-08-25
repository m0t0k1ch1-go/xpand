package xpand_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/xpand"
)

func TestFile(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		type input struct {
			path string
			opts []xpand.Option
		}

		tcs := []struct {
			name string
			in   input
			want string
		}{
			{
				"with a resolver for an empty scheme",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "env:XPAND_TEST_FOO" }}"}`),
					opts: []xpand.Option{
						xpand.WithResolver("", xpand.ResolveRaw),
					},
				},
				"invalid option: scheme must not be empty",
			},
			{
				"with a resolver for a scheme containing the reference separator",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "env:XPAND_TEST_FOO" }}"}`),
					opts: []xpand.Option{
						xpand.WithResolver("foo:bar", xpand.ResolveRaw),
					},
				},
				`invalid option: scheme "foo:bar" must not contain reference separator ":"`,
			},
			{
				"with a nil resolver",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "env:XPAND_TEST_FOO" }}"}`),
					opts: []xpand.Option{
						xpand.WithResolver("nil", nil),
					},
				},
				`invalid option: resolver for scheme "nil" must not be nil`,
			},
			{
				"with an unresolvable reference",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "env:XPAND_TEST_UNSET" }}"}`),
				},
				`no value resolved for "env:XPAND_TEST_UNSET"`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				_, err := xpand.File(ctx, tc.in.path, tc.in.opts...)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		type input struct {
			path string
			opts []xpand.Option
		}

		t.Setenv("XPAND_TEST_FOO", "foo")
		t.Setenv("XPAND_TEST_JSON", `{"foo":"bar"}`)

		tcs := []struct {
			name string
			in   input
			want string
		}{
			{
				"lookup with default delimiters",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "env:XPAND_TEST_FOO" }}"}`),
				},
				`{"foo":"foo"}`,
			},
			{
				"lookup with custom delimiters",
				input{
					path: writeFile(t, "config.json", `{"foo":"<< lookup "env:XPAND_TEST_FOO" >>"}`),
					opts: []xpand.Option{
						xpand.WithDelims("<<", ">>"),
					},
				},
				`{"foo":"foo"}`,
			},
			{
				"lookup with empty delimiters",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "env:XPAND_TEST_FOO" }}"}`),
					opts: []xpand.Option{
						xpand.WithDelims("", ""),
					},
				},
				`{"foo":"foo"}`,
			},
			{
				"lookup with a custom resolver",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "secret:foo" }}"}`),
					opts: []xpand.Option{
						xpand.WithResolver("secret", func(_ context.Context, key string) (string, bool, error) {
							v, ok := map[string]string{
								"foo": "s3cr3t",
							}[key]

							return v, ok, nil
						}),
					},
				},
				`{"foo":"s3cr3t"}`,
			},
			{
				"lookup with a custom resolver overriding a built-in one",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "raw:foo" }}"}`),
					opts: []xpand.Option{
						xpand.WithResolver(xpand.SchemeRaw, func(_ context.Context, key string) (string, bool, error) {
							return key + ".overridden", true, nil
						}),
					},
				},
				`{"foo":"foo.overridden"}`,
			},
			{
				"lookup with jsonEscape",
				input{
					path: writeFile(t, "config.txt", `{{ lookup "env:XPAND_TEST_JSON" | jsonEscape }}`),
				},
				`{\"foo\":\"bar\"}`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				b, err := xpand.File(ctx, tc.in.path, tc.in.opts...)
				require.NoError(t, err)
				require.Equal(t, tc.want, string(b))
			})
		}
	})
}

func writeFile(t *testing.T, name string, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)

	err := os.WriteFile(path, []byte(content), 0o600)
	require.NoError(t, err)

	return path
}
