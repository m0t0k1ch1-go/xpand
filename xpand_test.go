package xpand_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/xpand"
)

func TestFile(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"template function returns an error",
				`{"baz":"{{ lookup "env:XPAND_TEST_BAZ" }}"}`,
				`no value resolved for "env:XPAND_TEST_BAZ"`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				path := writeFile(t, "config.json", tc.in)

				_, err := xpand.File(ctx, path)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		type input struct {
			path string
			opts []xpand.Option
		}

		for k, v := range map[string]string{
			"XPAND_TEST_FOO": "foo",
			"XPAND_TEST_BAR": "bar",
		} {
			t.Setenv(k, v)
		}

		tcs := []struct {
			name string
			in   input
			want string
		}{
			{
				"with default delimiters",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "env:XPAND_TEST_FOO" "raw:default" }}","bar":"{{ lookup "env:XPAND_TEST_BAR" }}"}`),
					opts: nil,
				},
				`{"foo":"foo","bar":"bar"}`,
			},
			{
				"with custom delimiters",
				input{
					path: writeFile(t, "config.json", `{"foo":"<< lookup "env:XPAND_TEST_FOO" "raw:default" >>","bar":"<< lookup "env:XPAND_TEST_BAR" >>"}`),
					opts: []xpand.Option{
						xpand.WithDelims("<<", ">>"),
					},
				},
				`{"foo":"foo","bar":"bar"}`,
			},
			{
				"with empty delimiters",
				input{
					path: writeFile(t, "config.json", `{"foo":"{{ lookup "env:XPAND_TEST_FOO" "raw:default" }}","bar":"{{ lookup "env:XPAND_TEST_BAR" }}"}`),
					opts: []xpand.Option{
						xpand.WithDelims("", ""),
					},
				},
				`{"foo":"foo","bar":"bar"}`,
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
