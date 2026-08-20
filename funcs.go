package xpand

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

func newFuncMap(ctx context.Context, resolverMap map[string]resolver) template.FuncMap {
	return template.FuncMap{
		"lookup": func(refs ...string) (string, error) {
			return lookup(ctx, resolverMap, refs...)
		},
	}
}

func lookup(ctx context.Context, resolverMap map[string]resolver, refs ...string) (string, error) {
	if len(refs) == 0 {
		return "", errors.New("at least one reference must be provided")
	}

	for _, ref := range refs {
		scheme, key, ok := strings.Cut(ref, ":")
		if !ok {
			return "", fmt.Errorf("missing scheme in reference %q", ref)
		}

		r, ok := resolverMap[scheme]
		if !ok {
			return "", fmt.Errorf("unknown scheme %q in reference %q", scheme, ref)
		}

		v, ok, err := r(ctx, key)
		if err != nil {
			return "", fmt.Errorf("failed to resolve %q: %w", ref, err)
		}
		if ok {
			return v, nil
		}
	}

	quotedRefs := make([]string, len(refs))
	{
		for idx, ref := range refs {
			quotedRefs[idx] = strconv.Quote(ref)
		}
	}

	return "", fmt.Errorf("no value resolved for %s", strings.Join(quotedRefs, " or "))
}
