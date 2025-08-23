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
	loadFA()
}

func loadEN() {
	scope := "telegram-messages.loadEN"
	root, err := utils.GetRootOfProject()
	if err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}

	path := filepath.Join(root, "telegram-messages-en.yaml")

	eng := koanf.New(".")
	if err := eng.Load(file.Provider(path), yaml.Parser()); err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}
	k.MergeAt(eng, "en")
}

func loadFA() {
	scope := "telegram-messages.loadFA"
	root, err := utils.GetRootOfProject()
	if err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}

	path := filepath.Join(root, "telegram-messages-fa.yaml")

	fa := koanf.New(".")
	if err := fa.Load(file.Provider(path), yaml.Parser()); err != nil {
		momoError.Wrap(err).Scope(scope).Fatal()
	}
	k.MergeAt(fa, "fa")
}
