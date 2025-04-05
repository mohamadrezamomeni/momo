package log

import (
	"fmt"
	"log"
	"os"
	"time"

	"momo/pkg/utils"
)

type ErrorLog struct {
	kind    Kind
	message string
}

type ILog interface {
	WriteWarrning(string)
	WriteInfo(string)
}

type Log struct {
	accessFile string
	errorFile  string
}

func New(cfg LogConfig) ILog {
	if len(cfg.AccessFile) != 0 {
		utils.CreateFile(cfg.AccessFile)
	}
	if len(cfg.ErrorFile) != 0 {
		utils.CreateFile(cfg.ErrorFile)
	}
	return &Log{
		accessFile: cfg.AccessFile,
		errorFile:  cfg.ErrorFile,
	}
}

func (l *Log) getRightFile(kind Kind) (string, error) {
	errors := []string{warrning}
	acceesses := []string{info}
	for _, name := range errors {
		if name == kind {
			return l.errorFile, nil
		}
	}

	for _, name := range acceesses {
		if name == kind {
			return l.accessFile, nil
		}
	}

	return "", fmt.Errorf("kind doesnt match in system")
}

func (l *Log) logger(errorLog *ErrorLog) {
	now := time.Now().Format(time.RFC3339)

	row := fmt.Sprintf("[%s] %s %s", now, errorLog.kind, errorLog.message)

	fileName, err := l.getRightFile(errorLog.kind)
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

func (l *Log) WriteWarrning(message string) {
	log := ErrorLog{
		"warrning",
		message,
	}
	l.logger(&log)
}

func (l *Log) WriteInfo(message string) {
	log := ErrorLog{
		"info",
		message,
	}
	l.logger(&log)
}
