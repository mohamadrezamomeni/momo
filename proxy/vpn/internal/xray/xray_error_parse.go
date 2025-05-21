package xray

import (
	"fmt"
	"strings"
)

func isNotFoundError(err error, tag string) bool {
	if strings.Contains(err.Error(), fmt.Sprintf("not found: %s", tag)) {
		return true
	}
	return false
}

func isExisted(err error, tag string) bool {
	if strings.Contains(err.Error(), fmt.Sprintf("existing tag found: %s", tag)) {
		return true
	}
	return false
}
