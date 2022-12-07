package server

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/Bornholm/deformd/internal/config"
	"github.com/Bornholm/deformd/internal/handler"
	"github.com/Bornholm/deformd/internal/handler/module"
	"github.com/pkg/errors"
	"github.com/wneessen/go-mail"
)

func (s *Server) getRequestContextHandler(ctx context.Context) *handler.Handler {
	formConfig := s.getRequestContextFormConfig(ctx)
	if formConfig == nil {
		return nil
	}

	handlerConfig := formConfig.Handler.Config

	options := make([]handler.OptionFunc, 0)

	maxDuration := 10 * time.Second
	if handlerConfig.MaxDuration.Seconds() != 0 {
		maxDuration = handlerConfig.MaxDuration
	}

	options = append(options, handler.WithMaxDuration(maxDuration))

	// Configure enabled modules
	modules := configureModules(formConfig.Handler.Config.Modules)

	options = append(options, handler.WithModules(modules...))

	handler := handler.New(string(formConfig.Handler.Script), options...)

	return handler
}

func configureModules(conf config.ModulesConfig) []handler.ModuleFactory {
	modules := make([]handler.ModuleFactory, 0)

	if conf.Email != nil {
		modules = append(modules, configureEmailModule(conf.Email))
	}

	if conf.Params != nil {
		modules = append(modules, configureParamsModule(conf.Params))
	}

	modules = append(
		modules,
		module.ConsoleModuleFactory(nil),
		module.MessageModuleFactory(),
		module.RedirectModuleFactory(),
	)

	return modules
}

type emailModuleConfig struct {
	Host     string             `mapstructure:"host"`
	Port     *int               `mapstructure:"port"`
	Username *string            `mapstructure:"username"`
	Password *string            `mapstructure:"password"`
	AuthType *mail.SMTPAuthType `mapstructure:"authType"`
}

func configureEmailModule(conf *config.EmailModuleConfig) handler.ModuleFactory {
	options := []mail.Option{}

	if conf.Port != nil {
		options = append(options, mail.WithPort(int(*conf.Port)))
	}

	if conf.Username != nil {
		options = append(options, mail.WithUsername(string(*conf.Username)))
	}

	if conf.Password != nil {
		options = append(options, mail.WithPassword(string(*conf.Password)))
	}

	if conf.AuthType != nil {
		options = append(options, mail.WithSMTPAuth(mail.SMTPAuthType(*conf.AuthType)))
	}

	if conf.InsecureSkipVerify != nil {
		options = append(options, mail.WithTLSConfig(&tls.Config{
			InsecureSkipVerify: bool(*conf.InsecureSkipVerify),
		}))
	}

	if conf.UseSSL != nil && *conf.UseSSL == true {
		options = append(options, mail.WithSSL())
	}

	if conf.TLSPolicy != nil {
		options = append(options, mail.WithTLSPolicy(mail.TLSPolicy(*conf.TLSPolicy)))
	}

	return module.EmailModuleFactory(
		string(conf.Host),
		options...,
	)
}

func configureParamsModule(conf *config.ParamsConfig) handler.ModuleFactory {
	var values map[string]interface{}

	if conf != nil {
		values = *conf
	} else {
		values = make(map[string]interface{})
	}

	return module.ParamsModuleFactory(values)
}

func moduleInstantiationErrorFactory(moduleName string, err error) handler.ModuleFactory {
	return func() (handler.Module, error) {
		return nil, errors.Wrapf(errors.WithStack(err), "could not instantiate module '%s'", moduleName)
	}
}
