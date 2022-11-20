package module

import (
	"fmt"

	"github.com/Bornholm/deformd/internal/handler"
	"github.com/dop251/goja"
	"github.com/pkg/errors"
)

const ConsoleModuleName = "console"

// ConsoleModule provides debug utilities.
type ConsoleModule struct {
	printf LogFunc
}

func (m *ConsoleModule) Name() string {
	return ConsoleModuleName
}

func (m *ConsoleModule) toArguments(args []goja.Value) []any {
	anyArgs := make([]any, 0)

	for _, a := range args {
		anyArgs = append(anyArgs, a.Export())
	}

	return anyArgs
}

func (m *ConsoleModule) error(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	ctx := assertContext(call.Argument(0), rt)
	args := m.toArguments(call.Arguments[1:])

	if len(args) == 0 {
		return nil
	}

	m.printf(ctx, fmt.Sprintf("[HANDLER][ERROR] %s", args[0]), args[1:]...)

	return nil
}

func (m *ConsoleModule) log(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	ctx := assertContext(call.Argument(0), rt)
	args := m.toArguments(call.Arguments[1:])

	m.printf(ctx, fmt.Sprintf("[HANDLER][LOG] %s", args[0]), args[1:]...)

	if len(args) == 0 {
		return nil
	}

	return nil
}

func (m *ConsoleModule) Export(export *goja.Object) {
	if err := export.Set("log", m.log); err != nil {
		panic(errors.Wrap(err, "could not set 'log' function"))
	}

	if err := export.Set("error", m.error); err != nil {
		panic(errors.Wrap(err, "could not set 'error' function"))
	}
}

func ConsoleModuleFactory(logFunc LogFunc) handler.ModuleFactory {
	if logFunc == nil {
		logFunc = defaultLogger
	}

	return func() (handler.Module, error) {
		return &ConsoleModule{
			printf: logFunc,
		}, nil
	}
}
