package config

import (
	"time"
)

type Handler struct {
	Script InterpolatedString `yaml:"script"`
	Config HandlerConfig      `yaml:"config"`
}

type HandlerConfig struct {
	MaxDuration time.Duration `yaml:"maxDuration"`
	Modules     ModulesConfig `yaml:"modules"`
}

type ModulesConfig struct {
	Email  *EmailModuleConfig `yaml:"email"`
	Params *ParamsConfig      `yaml:"params"`
}

type EmailModuleConfig struct {
	Host               InterpolatedString  `yaml:"host"`
	Port               *InterpolatedInt    `yaml:"port"`
	Username           *InterpolatedString `yaml:"username"`
	Password           *InterpolatedString `yaml:"password"`
	AuthType           *InterpolatedString `yaml:"authType"`
	InsecureSkipVerify *InterpolatedBool   `yaml:"insecureSkipVerify"`
	TLSPolicy          *InterpolatedInt    `yaml:"tlsPolicy"`
	UseSSL             *InterpolatedBool   `yaml:"useSSL"`
}

type ParamsConfig = InterpolatedMap
