package xpand

import (
	"context"
	"encoding/json/jsontext"
	"errors"
	"fmt"
	"strconv"
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
		return "", errors.New("at least one reference must be provided")
	}

	for _, ref := range refs {
		scheme, key, ok := strings.Cut(ref, referenceSeparator)
		if !ok {
			return "", fmt.Errorf("missing scheme in reference %q", ref)
		}

		resolver, ok := resolverMap[scheme]
		if !ok {
			return "", fmt.Errorf("unknown scheme %q in reference %q", scheme, ref)
		}

		v, ok, err := resolver(ctx, key)
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

func jsonEscape(s string) (string, error) {
	b, err := jsontext.AppendQuote(nil, s)
	if err != nil {
		return "", fmt.Errorf("failed to append quote: %w", err)
	}

	return string(b[1 : len(b)-1]), nil
}
