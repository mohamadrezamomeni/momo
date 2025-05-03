package log

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mohamadrezamomeni/momo/pkg/utils"
)

type Record struct {
	kind    Kind
	message string
}

var (
	accessFile string
	errorFile  string
)

func Init(cfg LogConfig) {
	if len(cfg.AccessFile) != 0 {
		utils.CreateFile(cfg.AccessFile)
		accessFile = cfg.AccessFile
	}
	if len(cfg.ErrorFile) != 0 {
		utils.CreateFile(cfg.ErrorFile)
		errorFile = cfg.ErrorFile
	}
}

func getRightFile(kind Kind) (string, error) {
	errors := []string{warrning, debugging}
	acceesses := []string{info}
	for _, name := range errors {
		if name == kind {
			return errorFile, nil
		}
	}

	for _, name := range acceesses {
		if name == kind {
			return accessFile, nil
		}
	}

	return "", fmt.Errorf("kind doesnt match in system")
}

func logger(record *Record) {
	now := time.Now().Format(time.RFC3339)

	row := fmt.Sprintf("[%s] %s %s", now, record.kind, record.message+"\n")

	fileName, err := getRightFile(record.kind)
	if err != nil {
		log.Fatal("project has a ciritical, bug please report the issue")
	}

	if len(fileName) == 0 {
		return
	}
	file, err := os.OpenFile(fileName, os.O_APPEND, 0o644)
	if err != nil {
		return
	}

	defer file.Close()
	fmt.Println(row)
	_, _ = file.WriteString(row)
}

func Warrningf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	log := Record{
		"warrning",
		s,
	}
	logger(&log)
}

func Warrning(s string) {
	log := Record{
		"warrning",
		s,
	}
	logger(&log)
}

func Infof(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	log := Record{
		"info",
		s,
	}
	logger(&log)
}

func Info(s string) {
	log := Record{
		"info",
		s,
	}
	logger(&log)
}

func Debugging(s string) {
	log := Record{
		"debugging",
		s,
	}
	logger(&log)
}

func Debuggingf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	log := Record{
		"debugging",
		s,
	}
	logger(&log)
}
