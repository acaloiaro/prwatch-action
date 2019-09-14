package internal

import (
	"os"
	"testing"
)

func TestServiceInitialization(t *testing.T) {

	os.Setenv("JIRA_USER", "foo")
	os.Setenv("JIRA_API_TOKEN", "bar")
	os.Setenv("JIRA_HOST", "test.dev")

	if services.issues() == nil {
		t.Error("services provider should initialize its issue provider")
	}

	if services.git() == nil {
		t.Error("services provider should initialize its git provider")
	}

	if services.files() == nil {
		t.Error("services provider should initialize its files provider")
	}

}
