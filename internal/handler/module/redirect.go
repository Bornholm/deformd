package module

import (
	"context"

	"github.com/Bornholm/deformd/internal/handler"
	"github.com/dop251/goja"
	"github.com/pkg/errors"
)

const RedirectModuleName = "redirect"

// MessageModule provides redirection utilities.
type RedirectModule struct{}

func (m *RedirectModule) Name() string {
	return RedirectModuleName
}

func (m *RedirectModule) Export(export *goja.Object) {
	if err := export.Set("to", m.to); err != nil {
		panic(errors.Wrap(err, "could not set 'success' function"))
	}
}

func (m *RedirectModule) to(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	ctx := assertContext(call.Argument(0), rt)

	url, ok := call.Argument(1).Export().(string)
	if !ok {
		panic(errors.New("second argument should be a string"))
	}

	message, ok := call.Argument(2).Export().(string)
	if ok {
		if err := SetRedirectMessage(ctx, message); err != nil {
			panic(errors.Wrap(err, "could not set redirect message on context"))
		}
	}

	if err := SetRedirectURL(ctx, url); err != nil {
		panic(errors.Wrap(err, "could not set redirect url on context"))
	}

	return nil
}

func RedirectModuleFactory() handler.ModuleFactory {
	return func() (handler.Module, error) {
		return &RedirectModule{}, nil
	}
}

const redirectURLContextKey contextKey = "redirectURL"

func WithRedirectURL(ctx context.Context) (*string, context.Context) {
	redirectURL := ""
	ctx = context.WithValue(ctx, redirectURLContextKey, &redirectURL)

	return &redirectURL, ctx
}

func SetRedirectURL(ctx context.Context, url string) error {
	urlPtr, err := GetRedirectURL(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	*urlPtr = url

	return nil
}

func GetRedirectURL(ctx context.Context) (*string, error) {
	redirectURL, ok := ctx.Value(redirectURLContextKey).(*string)
	if !ok {
		return nil, errors.New("could not retrieve redirect url on context")
	}

	return redirectURL, nil
}

const redirectMessageContextKey contextKey = "redirectMessage"

func GetRedirectMessage(ctx context.Context) (*string, error) {
	redirectMessage, ok := ctx.Value(redirectMessageContextKey).(*string)
	if !ok {
		return nil, errors.New("could not retrieve redirect message on context")
	}

	return redirectMessage, nil
}

func WithRedirectMessage(ctx context.Context) (*string, context.Context) {
	redirectMessage := "You will be redirected in {{ .Delay}}s to '{{ .URL }}'..."
	ctx = context.WithValue(ctx, redirectMessageContextKey, &redirectMessage)

	return &redirectMessage, ctx
}

func SetRedirectMessage(ctx context.Context, message string) error {
	redirectMessagePtr, err := GetRedirectMessage(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	*redirectMessagePtr = message

	return nil
}
