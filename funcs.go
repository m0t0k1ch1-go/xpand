package xpand

import (
	"context"
	"encoding/json/jsontext"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

const (
	referenceSeparator = ":"
)

func newFuncMap(ctx context.Context, resolverMap map[string]Resolver) template.FuncMap {
	return template.FuncMap{
		"lookup": func(refs ...string) (string, error) {
			return lookup(ctx, resolverMap, refs...)
		},
		"jsonEscape": jsonEscape,
	}
}

func lookup(ctx context.Context, resolverMap map[string]Resolver, refs ...string) (string, error) {
	if len(refs) == 0 {
		return "", errors.New("invalid references: empty")
	}

	for _, ref := range refs {
		scheme, key, ok := strings.Cut(ref, referenceSeparator)
		if !ok {
			return "", fmt.Errorf("invalid reference %q: missing scheme", ref)
		}

		resolver, ok := resolverMap[scheme]
		if !ok {
			return "", fmt.Errorf("unsupported scheme: %q", scheme)
		}

		v, ok, err := resolver(ctx, key)
		if err != nil {
			return "", fmt.Errorf("failed to resolve %q: %w", ref, err)
		}
		if ok {
			return v, nil
		}
	}

	return "", errors.New("invalid references: no value resolved")
}

func jsonEscape(s string) (string, error) {
	b, err := jsontext.AppendQuote(nil, s)
	if err != nil {
		return "", fmt.Errorf("failed to append quote: %w", err)
	}

	return string(b[1 : len(b)-1]), nil
}
