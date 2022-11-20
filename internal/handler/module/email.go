package module

import (
	"context"

	"github.com/Bornholm/deformd/internal/handler"
	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"github.com/wneessen/go-mail"
)

const EmailModuleName = "email"

// EmailModule provides email sending utilities.
type EmailModule struct {
	ctx     context.Context
	host    string
	options []mail.Option
}

func (m *EmailModule) Name() string {
	return EmailModuleName
}

func (m *EmailModule) Export(export *goja.Object) {
	if err := export.Set("send", m.send); err != nil {
		panic(errors.Wrap(err, "could not set 'send' function"))
	}
}

func (m *EmailModule) send(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	ctx := assertContext(call.Argument(0), rt)

	attrs, ok := call.Argument(1).Export().(map[string]interface{})
	if !ok {
		panic(errors.New("second argument should be an object"))
	}

	msg := mail.NewMsg()

	to, err := getStringSliceAttr(attrs, "to", false)
	if err != nil {
		panic(errors.WithStack(err))
	}

	if err := msg.To(to...); err != nil {
		panic(errors.WithStack(err))
	}

	cc, err := getStringSliceAttr(attrs, "cc", true)
	if err != nil {
		panic(errors.WithStack(err))
	}

	if err := msg.Cc(cc...); err != nil {
		panic(errors.WithStack(err))
	}

	bcc, err := getStringSliceAttr(attrs, "bcc", true)
	if err != nil {
		panic(errors.WithStack(err))
	}

	if err := msg.Bcc(bcc...); err != nil {
		panic(errors.WithStack(err))
	}

	from, err := getStringAttr(attrs, "from", false)
	if err != nil {
		panic(errors.WithStack(err))
	}

	if err := msg.From(from); err != nil {
		panic(errors.WithStack(err))
	}

	subject, err := getStringAttr(attrs, "subject", false)
	if err != nil {
		panic(errors.WithStack(err))
	}

	msg.Subject(subject)

	body, err := getStringAttr(attrs, "body", false)
	if err != nil {
		panic(errors.WithStack(err))
	}

	msg.SetBodyString(mail.TypeTextPlain, body)

	if err := m.sendMessage(ctx, msg); err != nil {
		panic(errors.WithStack(err))
	}

	return nil
}

func (m *EmailModule) sendMessage(ctx context.Context, msg *mail.Msg) error {
	client, err := mail.NewClient(
		m.host,
		m.options...,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if err := client.Close(); err != nil && !errors.Is(err, mail.ErrNoActiveConnection) {
			panic(errors.WithStack(err))
		}
	}()

	if err := client.DialAndSendWithContext(ctx, msg); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func getStringAttr(attrs map[string]interface{}, attrName string, optional bool) (string, error) {
	rawAttr, exists := attrs[attrName]
	if !exists {
		if optional {
			return "", nil
		}

		return "", errors.Errorf("could not find attribute '%s'", attrName)
	}

	raw, ok := rawAttr.(string)
	if !ok {
		return "", errors.Errorf("could not cast attribute '%s' to string", attrName)
	}

	return raw, nil
}

func getStringSliceAttr(attrs map[string]interface{}, attrName string, optional bool) ([]string, error) {
	rawAttr, exists := attrs[attrName]
	if !exists {
		if optional {
			return []string{}, nil
		}

		return nil, errors.Errorf("could not find attribute '%s'", attrName)
	}

	raw, ok := rawAttr.([]interface{})
	if !ok {
		return nil, errors.Errorf("could not cast attribute '%s' to []interface", attrName)
	}

	strSlice, err := toStringSlice(raw)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return strSlice, nil
}

func EmailModuleFactory(host string, options ...mail.Option) handler.ModuleFactory {
	return func() (handler.Module, error) {
		return &EmailModule{
			host:    host,
			options: options,
		}, nil
	}
}

func toStringSlice(iSlice []interface{}) ([]string, error) {
	strSlice := make([]string, 0, len(iSlice))

	for index, item := range iSlice {
		str, ok := item.(string)
		if !ok {
			return nil, errors.Errorf("item #%d could not be casted to string", index)
		}

		strSlice = append(strSlice, str)
	}

	return strSlice, nil
}
