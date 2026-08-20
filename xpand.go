package xpand

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"
)

const (
	DefaultLeftDelim  = "{{"
	DefaultRightDelim = "}}"
)

type options struct {
	leftDelim   string
	rightDelim  string
	resolverMap map[string]Resolver
}

func (o *options) validate() error {
	errs := []error{}

	for scheme, resolver := range o.resolverMap {
		if len(scheme) == 0 {
			errs = append(errs, errors.New("scheme must not be empty"))
		} else if strings.Contains(scheme, referenceSeparator) {
			errs = append(errs, fmt.Errorf("scheme %q must not contain reference separator %q", scheme, referenceSeparator))
		} else if resolver == nil {
			errs = append(errs, fmt.Errorf("resolver for scheme %q must not be nil", scheme))
		}
	}

	return errors.Join(errs...)
}

// Option configures the expansion.
type Option func(*options)

// WithDelims sets the action delimiters. Empty values stand for the defaults.
func WithDelims(left, right string) Option {
	return func(o *options) {
		o.leftDelim = left
		o.rightDelim = right
	}
}

// WithResolver registers the resolver for the scheme.
// It overrides the resolver already registered for the same scheme, including the built-in ones.
func WithResolver(scheme string, resolver Resolver) Option {
	return func(o *options) {
		o.resolverMap[scheme] = resolver
	}
}

// File reads the file at path, expands it as a [text/template], and returns the result.
func File(ctx context.Context, path string, opts ...Option) ([]byte, error) {
	o := &options{
		leftDelim:   DefaultLeftDelim,
		rightDelim:  DefaultRightDelim,
		resolverMap: newResolverMap(),
	}
	{
		for _, opt := range opts {
			opt(o)
		}
	}
	if err := o.validate(); err != nil {
		return nil, fmt.Errorf("invalid option: %w", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	tmpl, err := template.
		New(path).
		Delims(o.leftDelim, o.rightDelim).
		Funcs(newFuncMap(ctx, o.resolverMap)).
		Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	var buf bytes.Buffer
	{
		if err := tmpl.Execute(&buf, nil); err != nil {
			return nil, fmt.Errorf("failed to execute: %w", err)
		}
	}

	return buf.Bytes(), nil
}
