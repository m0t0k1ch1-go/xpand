package xpand_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/xpand"
)

func TestEnv(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []string
			want string
		}{
			{
				"no keys and no default value",
				[]string{},
				"at least one key and a default value must be provided",
			},
			{
				"one key and no default value",
				[]string{"XPAND_TEST_FOO"},
				"at least one key and a default value must be provided",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				_, err := xpand.Env(ctx, tc.in...)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		for k, v := range map[string]string{
			"XPAND_TEST_FOO":   "foo",
			"XPAND_TEST_BAR":   "bar",
			"XPAND_TEST_EMPTY": "",
		} {
			t.Setenv(k, v)
		}

		tcs := []struct {
			name string
			in   []string
			want string
		}{
			{
				"one key and a default value: first key is set",
				[]string{"XPAND_TEST_FOO", "default"},
				"foo",
			},
			{
				"one key and a default value: first key is set to an empty string",
				[]string{"XPAND_TEST_EMPTY", "default"},
				"",
			},
			{
				"one key and a default value: first key is not set",
				[]string{"XPAND_TEST_BAZ", "default"},
				"default",
			},
			{
				"two keys and a default value: first key is set",
				[]string{"XPAND_TEST_FOO", "XPAND_TEST_BAR", "default"},
				"foo",
			},
			{
				"two keys and a default value: first key is not set, second key is set",
				[]string{"XPAND_TEST_BAZ", "XPAND_TEST_BAR", "default"},
				"bar",
			},
			{
				"two keys and a default value: both keys are not set",
				[]string{"XPAND_TEST_BAZ", "XPAND_TEST_QUX", "default"},
				"default",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				got, err := xpand.Env(ctx, tc.in...)
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			})
		}
	})
}

func TestMustEnv(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []string
			want string
		}{
			{
				"no keys",
				[]string{},
				"at least one key must be provided",
			},
			{
				"one key: key is not set",
				[]string{"XPAND_TEST_BAZ"},
				`"XPAND_TEST_BAZ" must be set`,
			},
			{
				"two keys: both keys are not set",
				[]string{"XPAND_TEST_BAZ", "XPAND_TEST_QUX"},
				`"XPAND_TEST_BAZ" or "XPAND_TEST_QUX" must be set`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				_, err := xpand.MustEnv(ctx, tc.in...)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		for k, v := range map[string]string{
			"XPAND_TEST_FOO":   "foo",
			"XPAND_TEST_BAR":   "bar",
			"XPAND_TEST_EMPTY": "",
		} {
			t.Setenv(k, v)
		}

		tcs := []struct {
			name string
			in   []string
			want string
		}{
			{
				"one key: key is set",
				[]string{"XPAND_TEST_FOO"},
				"foo",
			},
			{
				"one key: key is set to an empty string",
				[]string{"XPAND_TEST_EMPTY"},
				"",
			},
			{
				"two keys: first key is set",
				[]string{"XPAND_TEST_FOO", "XPAND_TEST_BAR"},
				"foo",
			},
			{
				"two keys: first key is not set, second key is set",
				[]string{"XPAND_TEST_BAZ", "XPAND_TEST_BAR"},
				"bar",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				got, err := xpand.MustEnv(ctx, tc.in...)
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			})
		}
	})
}
