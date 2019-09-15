package internal

import (
	"testing"

	"github.com/acaloiaro/prwatch/internal/config"
)

func TestServiceInitialization(t *testing.T) {

	config.GlobalEnable(config.Jira)
	config.GlobalSet(config.JiraUser, "foo")
	config.GlobalSet(config.JiraHost, "host.dev")

	config.SetEnv("JIRA_API_TOKEN", "bar")

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
