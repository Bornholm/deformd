package config

type HTTPConfig struct {
	Host    InterpolatedString `yaml:"host"`
	Port    uint               `yaml:"port"`
	BaseURL InterpolatedString `yaml:"baseUrl"`
}

func NewDefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Host:    "0.0.0.0",
		Port:    3000,
		BaseURL: "/",
	}
}
