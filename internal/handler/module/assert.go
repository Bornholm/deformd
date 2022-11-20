package module

import (
	"context"

	"github.com/dop251/goja"
)

func assertContext(v goja.Value, r *goja.Runtime) context.Context {
	if c, ok := v.Export().(context.Context); ok {
		return c
	}

	panic(r.NewTypeError("value should be a context"))
}
