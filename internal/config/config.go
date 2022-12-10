package config

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config definition
type Config struct {
	HTTP    HTTPConfig            `yaml:"http"`
	Logger  LoggerConfig          `yaml:"logger"`
	Forms   map[string]FormConfig `yaml:"forms"`
	Include []InterpolatedString  `yaml:"include"`
}

func (c *Config) LoadIncludes(baseDir string) error {
	for _, inc := range c.Include {
		pattern := filepath.Join(baseDir, string(inc))

		matches, err := filepath.Glob(pattern)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, m := range matches {
			data, err := ioutil.ReadFile(m)
			if err != nil {
				return errors.Wrapf(err, "could not read file '%s'", m)
			}

			if err := yaml.Unmarshal(data, c); err != nil {
				return errors.Wrapf(err, "could not unmarshal configuration")
			}
		}
	}

	return nil
}

// NewFromFile retrieves the configuration from the given file
func NewFromFile(path string) (*Config, error) {
	config := NewDefault()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file '%s'", path)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal configuration")
	}

	baseDir := filepath.Dir(path)

	if err := config.LoadIncludes(baseDir); err != nil {
		return nil, errors.WithStack(err)
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
