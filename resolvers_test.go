package xpand_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/xpand"
)

func TestResolveRaw(t *testing.T) {
	type output struct {
		value string
		ok    bool
	}

	tcs := []struct {
		name string
		in   string
		want output
	}{
		{
			"key",
			"foo",
			output{"foo", true},
		},
		{
			"empty key",
			"",
			output{"", true},
		},
		{
			"key containing the separator",
			"https://example.com",
			output{"https://example.com", true},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			v, ok, err := xpand.ResolveRaw(ctx, tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.want.value, v)
			require.Equal(t, tc.want.ok, ok)
		})
	}
}

func TestResolveEnv(t *testing.T) {
	type output struct {
		value string
		ok    bool
	}

	for k, v := range map[string]string{
		"XPAND_TEST_FOO":   "foo",
		"XPAND_TEST_EMPTY": "",
	} {
		t.Setenv(k, v)
	}

	tcs := []struct {
		name string
		in   string
		want output
	}{
		{
			"key is set",
			"XPAND_TEST_FOO",
			output{"foo", true},
		},
		{
			"key is set to an empty string",
			"XPAND_TEST_EMPTY",
			output{"", true},
		},
		{
			"key is not set",
			"XPAND_TEST_BAZ",
			output{"", false},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			v, ok, err := xpand.ResolveEnv(ctx, tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.want.value, v)
			require.Equal(t, tc.want.ok, ok)
		})
	}
}
