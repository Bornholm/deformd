package handler

import (
	"github.com/dop251/goja"
)

type Module interface {
	Name() string
	Export(export *goja.Object)
}

type ModuleFactory func() (Module, error)
