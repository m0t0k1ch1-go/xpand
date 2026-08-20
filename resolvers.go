package xpand

import (
	"context"
	"os"
)

const (
	SchemeRaw = "raw"
	SchemeEnv = "env"
)

// Resolver resolves the value for key.
// The second return value reports whether the value was resolved.
type Resolver func(ctx context.Context, key string) (string, bool, error)

func newResolverMap() map[string]Resolver {
	return map[string]Resolver{
		SchemeRaw: resolveRaw,
		SchemeEnv: resolveEnv,
	}
}

func resolveRaw(_ context.Context, key string) (string, bool, error) {
	return key, true, nil
}

func resolveEnv(_ context.Context, key string) (string, bool, error) {
	v, ok := os.LookupEnv(key)

	return v, ok, nil
}
