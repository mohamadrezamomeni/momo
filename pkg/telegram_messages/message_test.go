package telegrammessages

import (
	"os"
	"testing"

	"github.com/mohamadrezamomeni/momo/entity"
)

func TestMain(m *testing.M) {
	Load()
	code := m.Run()
	os.Exit(code)
}

func TestMessageEN(t *testing.T) {
	msg, err := GetMessage("testing.registeration_testing", map[string]string{
		"username": "max",
	}, entity.EN)
	if err != nil {
		t.Errorf("someting went wrong the problem was %v", err)
	}
	if msg != "hello max" {
		t.Errorf(`we expected we would get hello max but we got %s`, msg)
	}

	_, err = GetMessage("testing.registeration_testingg", map[string]string{
		"username": "max",
	}, entity.EN)
	if err == nil {
		t.Error("we expected an error")
	}
}

func TestMessageFN(t *testing.T) {
	msg, err := GetMessage("testing.registeration_testing", map[string]string{
		"username": "max",
	}, entity.FA)
	if err != nil {
		t.Errorf("someting went wrong the problem was %v", err)
	}
	if msg != "سلام max" {
		t.Errorf(`we expected we would get hello max but we got %s`, msg)
	}

	_, err = GetMessage("testing.registeration_testingg", map[string]string{
		"username": "max",
	}, entity.FA)
	if err == nil {
		t.Error("we expected an error")
	}
}
