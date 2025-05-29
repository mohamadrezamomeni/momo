package templates

import "testing"

func TestClientConfig(t *testing.T) {
	_, err := LoadClientConfig("instagram.com", "1234", "1234")
	if err != nil {
		t.Fatalf("someting went wrong that was %v", err)
	}
}
