package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitialize(t *testing.T) {

	configPath := "./github-actions/prwatch-action/config.yaml"
	defer os.Remove(configPath)

	writeConfig(configPath)

	Initialize()

	// global settings

	if SettingEnabled(IssueComments) {
		t.Errorf("global setting '%s' should be disabled", IssueComments)
	}

	if !SettingEnabled(IssueTransitions) {
		t.Errorf("global setting '%s' should be enabled", IssueTransitions)
	}

	if UserSettingEnabled("acaloiaro", IssueComments) {
		t.Errorf("setting should be disabled for acaloiaro: %s", IssueComments)
	}

	if !UserSettingEnabled("acaloiaro", IssueTransitions) {
		t.Errorf("setting should be enabled for acaloiaro: %s", IssueTransitions)
	}

	// user-specific settings

	Reset()

	UserEnable("foobar", IssueTransitions)
	if !UserSettingEnabled("foobar", IssueTransitions) {
		t.Error("user-specific setting should be enabled for foobar")
	}

	UserDisableSetting("foobar", IssueTransitions)
	if UserSettingEnabled("foobar", IssueTransitions) {
		t.Error("user-specific setting should be disabled for foobar")
	}

	UserEnable("foobar", IssueComments)
	if !UserSettingEnabled("foobar", IssueComments) {
		t.Error("user-specific setting should be enabled for foobar")
	}

	UserDisableSetting("foobar", IssueComments)
	if UserSettingEnabled("foobar", IssueComments) {
		t.Error("user-specific setting should be disabled for foobar")
	}

	GlobalDisable(IssueComments)
	UserEnable("foobar", IssueComments)
	if UserSettingEnabled("foobar", IssueComments) {
		t.Error("user-specific setting should be overridden by global setting")
	}

	Reset()

	GlobalDisable(IssueTransitions)
	UserEnable("foobar", IssueTransitions)
	if UserSettingEnabled("foobar", IssueTransitions) {
		t.Error("user-specific setting should be overridden by global setting")
	}

}

func writeConfig(path string) {
	var yaml = []byte(`---
settings:
  dual_pass:
    enabled: true
    wait_duration: 10s
  jira:
    enabled: true
    user: jira-bot
    host: greenhouseio.atlassian.net
    project_name: GREEN
  issues:
    conflict_status: In Progress
    enable_comment: false
    enable_transition: true
users:
  acaloiaro:
    settings:
      issues:
        enable_comment: true
        enable_transition: true
`)

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(yaml))
	if err != nil {
		fmt.Printf("Unable to parse config: %v", err)
	}

	folderPath := filepath.Dir(path)
	os.MkdirAll(folderPath, os.ModePerm)

	err = viper.WriteConfigAs(path)
	if err != nil {
		fmt.Printf("Unable to write config file: %v", err)
	}
}
