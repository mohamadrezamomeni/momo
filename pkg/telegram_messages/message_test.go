package telegrammessages

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Load()
	code := m.Run()
	os.Exit(code)
}

func TestMessage(t *testing.T) {
	msg, err := GetMessage("testing.registeration_testing", map[string]string{
		"username": "max",
	})
	if err != nil {
		t.Errorf("someting went wrong the problem was %v", err)
	}
	if msg != "hello max" {
		t.Errorf(`we expected we would get hello max but we got %s`, msg)
	}

	_, err = GetMessage("testing.registeration_testingg", map[string]string{
		"username": "max",
	})
	if err == nil {
		t.Error("we expected an error")
	}
}
