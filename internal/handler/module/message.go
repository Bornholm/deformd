package module

import (
	"context"

	"github.com/Bornholm/deformd/internal/handler"
	"github.com/dop251/goja"
	"github.com/pkg/errors"
)

const MessageModuleName = "message"

// MessageModulle provides message utilities.
type MessageModule struct{}

func (m *MessageModule) Name() string {
	return MessageModuleName
}

func (m *MessageModule) Export(export *goja.Object) {
	if err := export.Set("success", m.success); err != nil {
		panic(errors.Wrap(err, "could not set 'success' function"))
	}

	if err := export.Set("info", m.info); err != nil {
		panic(errors.Wrap(err, "could not set 'info' function"))
	}

	if err := export.Set("warn", m.warn); err != nil {
		panic(errors.Wrap(err, "could not set 'warn' function"))
	}

	if err := export.Set("error", m.error); err != nil {
		panic(errors.Wrap(err, "could not set 'error' function"))
	}
}

func (m *MessageModule) success(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	return m.message(MessageTypeSuccess, call, rt)
}

func (m *MessageModule) info(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	return m.message(MessageTypeInfo, call, rt)
}

func (m *MessageModule) warn(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	return m.message(MessageTypeWarn, call, rt)
}

func (m *MessageModule) error(call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	return m.message(MessageTypeError, call, rt)
}

func (m *MessageModule) message(messageType MessageType, call goja.FunctionCall, rt *goja.Runtime) goja.Value {
	ctx := assertContext(call.Argument(0), rt)

	text, ok := call.Argument(1).Export().(string)
	if !ok {
		panic(errors.New("second argument should be a string"))
	}

	messageStack, err := GetMessageStack(ctx)
	if err != nil {
		panic(errors.Wrap(err, "could not retrieve message stack on context"))
	}

	messageStack.Add(messageType, text)

	return nil
}

func MessageModuleFactory() handler.ModuleFactory {
	return func() (handler.Module, error) {
		return &MessageModule{}, nil
	}
}

type MessageStack struct {
	messages []*Message
}

func (s *MessageStack) Add(messageType MessageType, text string) {
	s.messages = append(s.messages, &Message{
		Type: messageType,
		Text: text,
	})
}

func (s *MessageStack) All() []*Message {
	return s.messages
}

func (s *MessageStack) HasError() bool {
	for _, m := range s.messages {
		if m.Type == MessageTypeError {
			return true
		}
	}

	return false
}

func NewMessageStack(messages ...*Message) *MessageStack {
	return &MessageStack{messages}
}

type MessageType int

const (
	MessageTypeSuccess MessageType = iota
	MessageTypeInfo
	MessageTypeWarn
	MessageTypeError
)

type Message struct {
	Type MessageType
	Text string
}

const messageContextKey contextKey = "message"

func WithNewMessageStack(ctx context.Context) (*MessageStack, context.Context) {
	messageStack := &MessageStack{
		messages: make([]*Message, 0),
	}
	ctx = context.WithValue(ctx, messageContextKey, messageStack)

	return messageStack, ctx
}

func GetMessageStack(ctx context.Context) (*MessageStack, error) {
	messageStack, ok := ctx.Value(messageContextKey).(*MessageStack)
	if !ok {
		return nil, errors.New("could not retrieve message stack on context")
	}

	return messageStack, nil
}
