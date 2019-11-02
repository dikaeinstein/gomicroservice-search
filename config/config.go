package config

import (
	"encoding/base64"

	"github.com/kelseyhightower/envconfig"
)

// Config represents env values
type Config struct {
	DogStatsD       string `envconfig:"DOGSTATSD"`
	MysqlConnection string `split_words:"true"`
	RsaPublicKey    string `split_words:"true"`
}

// New initializes and returns a new config
func New() Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	return cfg
}

// ParseRsaPublicKeyHex parses the RsaPublicKey hex and return bytes slice
func ParseRsaPublicKeyHex(k string) ([]byte, error) {
	rsaBytes, err := base64.StdEncoding.DecodeString(k)
	if err != nil {
		return nil, err
	}

	return rsaBytes, nil
}
