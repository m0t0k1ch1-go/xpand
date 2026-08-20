package xpand

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func newFuncMap(ctx context.Context) template.FuncMap {
	return template.FuncMap{
		"env": func(keysAndDefaultValue ...string) (string, error) {
			return env(ctx, keysAndDefaultValue...)
		},
		"must_env": func(keys ...string) (string, error) {
			return mustEnv(ctx, keys...)
		},
	}
}

func env(_ context.Context, keysAndDefaultValue ...string) (string, error) {
	if len(keysAndDefaultValue) < 2 {
		return "", errors.New("at least one key and a default value must be provided")
	}

	keys := keysAndDefaultValue[:len(keysAndDefaultValue)-1]
	defaultValue := keysAndDefaultValue[len(keysAndDefaultValue)-1]

	for _, key := range keys {
		if v, ok := os.LookupEnv(key); ok {
			return v, nil
		}
	}

	return defaultValue, nil
}

func mustEnv(_ context.Context, keys ...string) (string, error) {
	if len(keys) == 0 {
		return "", errors.New("at least one key must be provided")
	}

	for _, key := range keys {
		if v, ok := os.LookupEnv(key); ok {
			return v, nil
		}
	}

	quotedKeys := make([]string, len(keys))
	{
		for idx, key := range keys {
			quotedKeys[idx] = strconv.Quote(key)
		}
	}

	return "", fmt.Errorf("%s must be set", strings.Join(quotedKeys, " or "))
}
