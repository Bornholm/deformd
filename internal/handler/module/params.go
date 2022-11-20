package module

import (
	"github.com/Bornholm/deformd/internal/handler"
	"github.com/dop251/goja"
	"github.com/pkg/errors"
)

const ParamsModuleName = "params"

// ParamsModule provides parameters injection utilities.
type ParamsModule struct {
	values map[string]interface{}
}

func (m *ParamsModule) Name() string {
	return ParamsModuleName
}

func (m *ParamsModule) Export(export *goja.Object) {
	if err := export.Set("get", m.get); err != nil {
		panic(errors.Wrap(err, "could not set 'get' function"))
	}
}

func (m *ParamsModule) get(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	key, ok := call.Argument(0).Export().(string)
	if !ok {
		panic(errors.New("second argument should be an string"))
	}

	value, exists := m.values[key]
	if !exists {
		return rt.ToValue(nil)
	}

	return rt.ToValue(value)
}

func ParamsModuleFactory(values map[string]interface{}) handler.ModuleFactory {
	return func() (handler.Module, error) {
		return &ParamsModule{
			values: values,
		}, nil
	}
}
