package telegrammessages

import (
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

var k = koanf.New(".")

func Load() {
	loadEN()
}

func loadEN() {
	scope := "telegram-messages.loaden"
	root, err := utils.GetRootOfProject()
	if err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}

	path := filepath.Join(root, "telegram-messages-en.yaml")

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}
}
