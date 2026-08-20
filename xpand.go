package xpand

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"
)

const (
	DefaultLeftDelim  = "{{"
	DefaultRightDelim = "}}"
)

type options struct {
	leftDelim  string
	rightDelim string
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

// File reads the file at path, expands it as a [text/template], and returns the result.
func File(ctx context.Context, path string, opts ...Option) ([]byte, error) {
	o := &options{
		leftDelim:  DefaultLeftDelim,
		rightDelim: DefaultRightDelim,
	}
	for _, opt := range opts {
		opt(o)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	tmpl, err := template.New(path).Delims(o.leftDelim, o.rightDelim).Funcs(newFuncMap(ctx)).Parse(string(b))
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
