package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Bornholm/deformd/internal/config"
	"github.com/Bornholm/deformd/internal/form"
	"github.com/Bornholm/deformd/internal/handler/module"
	"github.com/Bornholm/deformd/internal/server/template"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/wpetit/goweb/logger"
)

type templateData struct {
	BaseURL  string
	Form     *form.Form
	Values   url.Values
	Messages *module.MessageStack
}

func (s *Server) serveForm(w http.ResponseWriter, r *http.Request) {
	form := s.getForm(w, r)
	if form == nil {
		return
	}

	messageStack, err := s.getFlashMessageStack(w, r)
	if err != nil {
		panic(errors.WithStack(err))
	}

	data := templateData{
		BaseURL:  string(s.conf.HTTP.BaseURL),
		Form:     form,
		Messages: messageStack,
	}

	if err := template.Exec("form.html.tmpl", w, data); err != nil {
		panic(errors.WithStack(err))
	}
}

func (s *Server) handleForm(w http.ResponseWriter, r *http.Request) {
	form := s.getForm(w, r)
	if form == nil {
		return
	}

	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		logger.Error(ctx, "could not parse form", logger.E(errors.WithStack(err)))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	handler := s.getRequestContextHandler(ctx)
	if handler == nil {
		logger.Error(ctx, "could not retrieve form handler")

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	messageStack, ctx := module.WithNewMessageStack(ctx)
	redirectURL, ctx := module.WithRedirectURL(ctx)

	if err := handler.Process(ctx, r.Form); err != nil {
		logger.Error(ctx, "could not process form", logger.E(errors.WithStack(err)))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	if messageStack.HasError() {
		data := templateData{
			BaseURL:  string(s.conf.HTTP.BaseURL),
			Form:     form,
			Values:   r.Form,
			Messages: messageStack,
		}

		if err := template.Exec("form.html.tmpl", w, data); err != nil {
			panic(errors.WithStack(err))
		}

		return
	}

	if err := s.setFlashMessageStack(w, messageStack); err != nil {
		panic(errors.WithStack(err))
	}

	if redirectURL != nil && *redirectURL != "" {
		if err := s.setFlashRedirectURL(w, *redirectURL); err != nil {
			panic(errors.WithStack(err))
		}

		http.Redirect(w, r, r.URL.String()+"/redirect", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
	}
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	form := s.getForm(w, r)
	if form == nil {
		return
	}

	messageStack, err := s.getFlashMessageStack(w, r)
	if err != nil {
		panic(errors.Wrap(err, "could not retrieve message stack"))
	}

	redirectURL, err := s.getFlashRedirectURL(w, r)
	if err != nil {
		panic(errors.Wrap(err, "could not retrieve redirect url"))
	}

	if err := template.Exec("redirect.html.tmpl", w, struct {
		BaseURL     string
		Messages    *module.MessageStack
		RedirectURL string
	}{
		BaseURL:     string(s.conf.HTTP.BaseURL),
		Messages:    messageStack,
		RedirectURL: redirectURL,
	}); err != nil {
		panic(errors.WithStack(err))
	}
}

func (s *Server) getRequestContextFormConfig(ctx context.Context) *config.FormConfig {
	formID := chi.URLParamFromCtx(ctx, "formID")

	logger.Debug(ctx, "retrieved request form id", logger.F("formID", formID))

	formConfig, exists := s.conf.Forms[formID]
	if !exists {
		return nil
	}

	return &formConfig
}

func (s *Server) getRequestContextForm(ctx context.Context) *form.Form {
	formConfig := s.getRequestContextFormConfig(ctx)
	if formConfig == nil {
		return nil
	}

	form := &form.Form{
		Title:  string(formConfig.Title),
		Fields: formConfig.Fields,
	}

	return form
}

func (s *Server) getForm(w http.ResponseWriter, r *http.Request) *form.Form {
	ctx := r.Context()
	form := s.getRequestContextForm(ctx)

	if form == nil {
		logger.Warn(ctx, "could not retrieve form from context", logger.F("url", r.RequestURI))

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return nil
	}

	return form
}
