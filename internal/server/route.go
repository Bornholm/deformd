package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Bornholm/deformd/internal/config"
	"github.com/Bornholm/deformd/internal/form"
	"github.com/Bornholm/deformd/internal/server/template"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gitlab.com/wpetit/goweb/logger"
)

type templateData struct {
	BaseURL string
	Form    *form.Form
	Values  url.Values
}

func (s *Server) serveForm(w http.ResponseWriter, r *http.Request) {
	form := s.getForm(w, r)
	if form == nil {
		return
	}

	data := templateData{
		BaseURL: string(s.conf.HTTP.BaseURL),
		Form:    form,
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

	if err := handler.Process(ctx, r.Form); err != nil {
		logger.Error(ctx, "could not process form", logger.E(errors.WithStack(err)))

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	data := templateData{
		BaseURL: string(s.conf.HTTP.BaseURL),
		Form:    form,
		Values:  r.Form,
	}

	if err := template.Exec("form.html.tmpl", w, data); err != nil {
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
