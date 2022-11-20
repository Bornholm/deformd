package handler

import "time"

type Options struct {
	Modules     []ModuleFactory
	MaxDuration time.Duration
}

type OptionFunc func(*Options)

func DefaultOptions() *Options {
	return &Options{
		Modules:     make([]ModuleFactory, 0),
		MaxDuration: time.Second * 10,
	}
}

func WithModules(modules ...ModuleFactory) OptionFunc {
	return func(options *Options) {
		options.Modules = modules
	}
}

func WithMaxDuration(duration time.Duration) OptionFunc {
	return func(options *Options) {
		options.MaxDuration = duration
	}
}
