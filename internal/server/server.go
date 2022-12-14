package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Bornholm/deformd/internal/config"
	"github.com/Bornholm/deformd/internal/server/template"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"gitlab.com/wpetit/goweb/logger"
)

type Server struct {
	conf  *config.Config
	store *sessions.CookieStore
}

type OnUpdateFunc func(values interface{}) error

func (s *Server) Start(ctx context.Context) (<-chan net.Addr, <-chan error) {
	errs := make(chan error)
	addrs := make(chan net.Addr)

	go s.run(ctx, addrs, errs)

	return addrs, errs
}

func (s *Server) run(parentCtx context.Context, addrs chan net.Addr, errs chan error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.conf.HTTP.Host, s.conf.HTTP.Port))
	if err != nil {
		errs <- errors.WithStack(err)

		return
	}

	addrs <- listener.Addr()

	defer func() {
		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			errs <- errors.WithStack(err)
		}

		close(errs)
		close(addrs)
	}()

	go func() {
		<-ctx.Done()

		if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("%+v", errors.WithStack(err))
		}
	}()

	templates := getEmbeddedTemplates()

	if err := template.Load(templates, "template"); err != nil {
		errs <- errors.WithStack(err)

		return
	}

	assets := getEmbeddedAssets()
	assetsHandler := http.FileServer(http.FS(assets))

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/forms/{formID}", s.serveForm)
	router.Post("/forms/{formID}", s.handleForm)
	router.Get("/forms/{formID}/redirect", s.handleRedirect)
	router.Handle("/assets/*", assetsHandler)

	logger.Info(ctx, "http server listening")

	if err := http.Serve(listener, router); err != nil && !errors.Is(err, net.ErrClosed) {
		errs <- errors.WithStack(err)
	}

	logger.Info(ctx, "http server exiting")
}

func New(funcs ...OptionFunc) *Server {
	opt := defaultOption()
	for _, fn := range funcs {
		fn(opt)
	}

	store := sessions.NewCookieStore()

	return &Server{
		conf:  opt.Config,
		store: store,
	}
}
