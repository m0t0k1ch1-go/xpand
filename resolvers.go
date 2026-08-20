package xpand

import (
	"context"
	"os"
)

type resolver func(ctx context.Context, key string) (string, bool, error)

func newResolverMap() map[string]resolver {
	return map[string]resolver{
		"raw": resolveRaw,
		"env": resolveEnv,
	}
}

func resolveRaw(_ context.Context, key string) (string, bool, error) {
	return key, true, nil
}

func resolveEnv(_ context.Context, key string) (string, bool, error) {
	v, ok := os.LookupEnv(key)

	return v, ok, nil
}
