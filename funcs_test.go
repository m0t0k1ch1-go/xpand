package xpand_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/xpand"
)

func TestLookup(t *testing.T) {
	resolverMap := map[string]xpand.Resolver{
		"is": func(_ context.Context, key string) (string, bool, error) {
			return key, true, nil
		},
		"empty": func(_ context.Context, _ string) (string, bool, error) {
			return "", true, nil
		},
		"skip": func(_ context.Context, _ string) (string, bool, error) {
			return "", false, nil
		},
		"fail": func(_ context.Context, _ string) (string, bool, error) {
			return "", false, errors.New("something went wrong")
		},
	}

	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []string
			want string
		}{
			{
				"no references",
				[]string{},
				"at least one reference must be provided",
			},
			{
				"single reference: missing scheme",
				[]string{"foo"},
				`missing scheme in reference "foo"`,
			},
			{
				"single reference: unknown scheme",
				[]string{"unknown:foo"},
				`unknown scheme "unknown" in reference "unknown:foo"`,
			},
			{
				"single reference: resolver returns an error",
				[]string{"fail:foo"},
				`failed to resolve "fail:foo": something went wrong`,
			},
			{
				"single reference: not resolved",
				[]string{"skip:foo"},
				`no value resolved for "skip:foo"`,
			},
			{
				"multiple references: none resolved",
				[]string{"skip:foo", "skip:bar"},
				`no value resolved for "skip:foo" or "skip:bar"`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				_, err := xpand.Lookup(ctx, resolverMap, tc.in...)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []string
			want string
		}{
			{
				"single reference: resolved",
				[]string{"is:foo"},
				"foo",
			},
			{
				"single reference: resolved to an empty value",
				[]string{"empty:foo"},
				"",
			},
			{
				"multiple references: the first is resolved",
				[]string{"is:foo", "fail:bar"},
				"foo",
			},
			{
				"multiple references: the first is not resolved, the second is",
				[]string{"skip:foo", "is:bar"},
				"bar",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				v, err := xpand.Lookup(ctx, resolverMap, tc.in...)
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			})
		}
	})
}

func TestJSONEscape(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   string
			want string
		}{
			{
				"empty",
				"",
				"",
			},
			{
				"nothing to escape",
				"foo",
				"foo",
			},
			{
				"tab and newline",
				"foo\tbar\n",
				`foo\tbar\n`,
			},
			{
				"json object",
				`{"foo":"bar"}`,
				`{\"foo\":\"bar\"}`,
			},
			{
				"html characters: not escaped",
				"<foo>&<bar>",
				"<foo>&<bar>",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := xpand.JSONEscape(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			})
		}
	})
}
