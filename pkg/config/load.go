package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	parser = yaml.Parser()
	k      = koanf.New(".")
)

func Load(path string) (Config, error) {
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("error loading config error: %v", err)
	}
	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config: %v", err)
	}
	return config, nil
}
