package server

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io"
	"net/http"
	"time"

	"github.com/Bornholm/deformd/internal/handler/module"
	"github.com/pkg/errors"
)

func init() {
	gob.Register(&gobMessageStack{})
}

type gobMessageStack struct {
	Messages []*module.Message
}

const (
	flashKeyMessageStack = "message-stack"
)

func (s *Server) getFlashMessageStack(w http.ResponseWriter, r *http.Request) (*module.MessageStack, error) {
	data, err := s.getFlash(w, r, flashKeyMessageStack)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	buf := bytes.NewBuffer(data)

	decoder := gob.NewDecoder(buf)

	messageStack := &gobMessageStack{}

	if err := decoder.Decode(messageStack); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}

		return nil, errors.WithStack(err)
	}

	return module.NewMessageStack(messageStack.Messages...), nil
}

func (s *Server) setFlashMessageStack(w http.ResponseWriter, stack *module.MessageStack) error {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(&gobMessageStack{stack.All()}); err != nil {
		return errors.WithStack(err)
	}

	s.setFlash(w, flashKeyMessageStack, buf.Bytes())

	return nil
}

func (s *Server) setFlash(w http.ResponseWriter, name string, value []byte) {
	cookie := &http.Cookie{Name: name, Value: encode(value)}

	http.SetCookie(w, cookie)
}

func (s *Server) getFlash(w http.ResponseWriter, r *http.Request, name string) ([]byte, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}

		return nil, errors.WithStack(err)
	}

	value, err := decode(cookie.Value)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dc := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}

	http.SetCookie(w, dc)

	return value, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
