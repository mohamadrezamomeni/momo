package templates

import (
	"fmt"
	"os"
	"path"
	"strings"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

func getFilename(name string) (string, error) {
	scope := "templates.getFilename"
	root, err := utils.GetRootOfProject()
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Input(name).ErrorWrite()
	}
	filename := fmt.Sprintf("%s.tpl", name)

	path := path.Join(root, "templates", filename)
	return path, nil
}

func LoadClientConfig(domain string, port string, userID string) (string, error) {
	scope := "template.LoadClientConfig"

	path, err := getFilename("client_config")
	if err != nil {
		return "", err
	}

	templateByte, err := os.ReadFile(path)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).Input(domain, port, userID).ErrorWrite()
	}
	template := string(templateByte)

	template = strings.ReplaceAll(template, fmt.Sprintf(":%s:", "domain"), domain)
	template = strings.ReplaceAll(template, fmt.Sprintf(":%s:", "port"), port)
	template = strings.ReplaceAll(template, fmt.Sprintf(":%s:", "user_id"), userID)

	return template, nil
}
