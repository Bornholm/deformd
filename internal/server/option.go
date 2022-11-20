package server

import "github.com/Bornholm/deformd/internal/config"

type Option struct {
	Config *config.Config
}

type OptionFunc func(*Option)

func defaultOption() *Option {
	return &Option{
		Config: config.NewDefault(),
	}
}

func WithConfig(conf *config.Config) OptionFunc {
	return func(opt *Option) {
		opt.Config = conf
	}
}
