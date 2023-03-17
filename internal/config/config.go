package config

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"gitlab.com/wpetit/goweb/logger"
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
	ctx := context.Background()

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, inc := range c.Include {
		pattern := filepath.Join(absBaseDir, string(inc))

		logger.Info(ctx, "searching inclusions", logger.F("pattern", pattern))

		matches, err := filepath.Glob(pattern)
		if err != nil {
			return errors.WithStack(err)
		}

		logger.Info(ctx, "included files", logger.F("matches", matches))

		for _, m := range matches {
			logger.Info(ctx, "loading included configuration", logger.F("file", m))

			data, err := ioutil.ReadFile(m)
			if err != nil {
				return errors.Wrapf(err, "could not read file '%s'", m)
			}

			if err := yaml.Unmarshal(data, c); err != nil {
				return errors.Wrapf(err, "could not unmarshal included configuration '%s'", m)
			}
		}
	}

	return nil
}

// NewFromFile retrieves the configuration from the given file
func NewFromFile(path string) (*Config, error) {
	config := NewDefault()

	logger.Info(context.Background(), "loading configuration", logger.F("path", path))

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

	logger.Debug(context.Background(), "loaded configuration", logger.F("config", config))

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
