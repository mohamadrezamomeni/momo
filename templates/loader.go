package templates

import (
	"fmt"
	"os"
	"strings"
)

func LoadClientConfig(domain string, port string, userID string) (string, error) {
	templateByte, err := os.ReadFile("client_config.tpl")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", err
	}
	template := string(templateByte)

	template = strings.ReplaceAll(template, fmt.Sprintf(":%s:", "domain"), domain)
	template = strings.ReplaceAll(template, fmt.Sprintf(":%s:", "port"), port)
	template = strings.ReplaceAll(template, fmt.Sprintf(":%s:", "user_id"), userID)

	return template, nil
}
