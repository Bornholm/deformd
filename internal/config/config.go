package config

import (
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config definition
type Config struct {
	HTTP   HTTPConfig            `yaml:"http"`
	Logger LoggerConfig          `yaml:"logger"`
	Forms  map[string]FormConfig `yaml:"forms"`
}

// NewFromFile retrieves the configuration from the given file
func NewFromFile(filepath string) (*Config, error) {
	config := NewDefault()

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file '%s'", filepath)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal configuration")
	}

	return config, nil
}

// NewDumpDefault dump the new default configuration
func NewDumpDefault() *Config {
	config := NewDefault()

	return config
}

// NewDefault return new default configuration
func NewDefault() *Config {
	return &Config{
		HTTP:   NewDefaultHTTPConfig(),
		Logger: NewDefaultLoggerConfig(),
	}
}

// Dump the given configuration in the given writer
func Dump(config *Config, w io.Writer) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "could not dump config")
	}

	if _, err := w.Write(data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
