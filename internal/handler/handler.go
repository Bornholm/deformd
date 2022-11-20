package handler

import (
	"context"
	"net/url"

	"github.com/dop251/goja"
	"github.com/pkg/errors"
)

const (
	RefHandleFunc = "handle"
)

type Handler struct {
	options *Options
	script  string
}

func (h *Handler) Process(ctx context.Context, form url.Values) error {
	// Create new ECMAscript runtime
	runtime := goja.New()

	runtime.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	// Load modules
	for _, factory := range h.options.Modules {
		mod, err := factory()
		if err != nil {
			return errors.Wrap(err, "could not instantiate module")
		}

		export := runtime.NewObject()
		mod.Export(export)

		if err := runtime.Set(mod.Name(), export); err != nil {
			return errors.Wrapf(err, "could not set module '%s'", mod.Name())
		}
	}

	var err error

	if err := runtime.Set("form", map[string][]string(form)); err != nil {
		return errors.WithStack(err)
	}

	ctx, cancel := context.WithTimeout(ctx, h.options.MaxDuration)
	defer cancel()

	if err := runtime.Set("ctx", ctx); err != nil {
		return errors.WithStack(err)
	}

	done := make(chan error)

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					done <- errors.WithStack(err)
				} else {
					panic(rec)
				}
			}

			close(done)
		}()

		_, err = runtime.RunString(h.script)
		if err != nil {
			done <- errors.WithStack(err)
		}

		done <- nil
	}()

	for {
		select {
		case <-ctx.Done():
			runtime.Interrupt(errors.WithStack(ctx.Err()))
		case err := <-done:
			if err != nil {
				return errors.WithStack(err)
			}

			return nil
		}
	}
}

func New(script string, funcs ...OptionFunc) *Handler {
	options := DefaultOptions()

	for _, fn := range funcs {
		fn(options)
	}

	return &Handler{
		options: options,
		script:  script,
	}
}
