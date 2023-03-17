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
	"gitlab.com/wpetit/goweb/logger"
)

func init() {
	gob.Register(&gobMessageStack{})
}

type gobMessageStack struct {
	Messages []*module.Message
}

const (
	flashKeyMessageStack    = "message-stack"
	flashKeyRedirectURL     = "redirect-url"
	flashKeyRedirectMessage = "redirect-message"
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

	s.clearFlash(w, flashKeyMessageStack)

	logger.Debug(r.Context(), "retrieved message stack", logger.F("messageStack", messageStack))

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

func (s *Server) getFlashRedirectURL(w http.ResponseWriter, r *http.Request) (string, error) {
	data, err := s.getFlash(w, r, flashKeyRedirectURL)
	if err != nil {
		return "", errors.WithStack(err)
	}

	s.clearFlash(w, flashKeyRedirectURL)

	redirectURL := string(data)

	logger.Debug(r.Context(), "retrieved redirect url", logger.F("redirectURL", redirectURL))

	return redirectURL, nil
}

func (s *Server) setFlashRedirectURL(w http.ResponseWriter, url string) error {
	s.setFlash(w, flashKeyRedirectURL, []byte(url))

	return nil
}

func (s *Server) getFlashRedirectMessage(w http.ResponseWriter, r *http.Request) (string, error) {
	data, err := s.getFlash(w, r, flashKeyRedirectMessage)
	if err != nil {
		return "", errors.WithStack(err)
	}

	s.clearFlash(w, flashKeyRedirectMessage)

	redirectMessage := string(data)

	logger.Debug(r.Context(), "retrieved redirect message", logger.F("redirectMessage", redirectMessage))

	return redirectMessage, nil
}

func (s *Server) setFlashRedirectMessage(w http.ResponseWriter, message string) error {
	s.setFlash(w, flashKeyRedirectMessage, []byte(message))

	return nil
}

func (s *Server) clearFlash(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{Name: name, Value: "", Path: "/"}

	http.SetCookie(w, cookie)
}

func (s *Server) setFlash(w http.ResponseWriter, name string, value []byte) {
	cookie := &http.Cookie{Name: name, Value: encode(value), Path: "/"}

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
