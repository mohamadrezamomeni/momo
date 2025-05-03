package config

import (
	"fmt"
	"path/filepath"

	"github.com/mohamadrezamomeni/momo/pkg/utils"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	parser = yaml.Parser()
	k      = koanf.New(".")
)

func Load(fileName string) (Config, error) {
	root, err := utils.GetRootOfProject()
	if err != nil {
		return Config{}, fmt.Errorf("somethine wrong happend error: %v", err)
	}
	path := filepath.Join(root, fileName)
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("error loading config error: %v", err)
	}
	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return Config{}, fmt.Errorf("error unmarshaling config: %v", err)
	}
	return config, nil
}
