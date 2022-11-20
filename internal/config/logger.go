package config

import "gitlab.com/wpetit/goweb/logger"

type LoggerConfig struct {
	Level  InterpolatedInt    `yaml:"level"`
	Format InterpolatedString `yaml:"format"`
}

func NewDefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:  InterpolatedInt(logger.LevelInfo),
		Format: InterpolatedString(logger.FormatHuman),
	}
}
