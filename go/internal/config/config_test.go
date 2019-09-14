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

	defer os.RemoveAll("./github-actions")

	// Make sure the config file can be read from various locations to simplify local dev
	configPath := "./config.yaml"
	writeConfig(configPath)

	Initialize()

	if !GetBool("settings.issue_transitions.enabled") {
		t.Error("configuration file should have enabled issue transitions")
	}

	os.Remove(configPath)
	Reset()

	configPath = "./github-actions/prwatch-action/config.yaml"
	writeConfig(configPath)

	Initialize()

	if !GetBool("settings.issue_transitions.enabled") {
		t.Error("configuration file should have enabled issue transitions")
	}
}

func TestSettingEnabled(t *testing.T) {
	configPath := "./config.yaml"
	writeConfig(configPath)

	Initialize()

	// global settings

	setting := "issue_comments"
	if SettingEnabled(setting) {
		t.Errorf("global setting '%s' should be disabled", setting)
	}

	setting = "issue_transitions"
	if !SettingEnabled(setting) {
		t.Errorf("global setting '%s' should be enabled", setting)
	}

	// user-specific settings

	setting = "issue_comments"
	if SettingEnabled("acaloiaro", setting) {
		t.Errorf("setting should be disabled for acaloiaro: %s", setting)
	}

	setting = "issue_transitions"
	if SettingEnabled("acaloiaro", setting) {
		t.Errorf("setting should be disabled for acaloiaro: %s", setting)
	}

	if SettingEnabled("foobar", "issue_transitions") {
		t.Error("user-specific setting should not be enabled for foobar")
	}

	UserEnable("foobar", "issue_transitions")
	if !SettingEnabled("foobar", "issue_transitions") {
		t.Error("user-specific setting should be enabled for foobar")
	}

	GlobalDisable("issue_transitions")
	if SettingEnabled("foobar", "issue_transitions") {
		t.Error("user-specific setting should be disabled for foobar")
	}
}

func writeConfig(path string) {
	var yaml = []byte(`---
settings:
  issue_transitions:
    enabled: true
  issue_comments:
    enabled: false
  dual_pass:
    enabled: true
    wait_duration: 1s

users:
  acaloiaro:
    settings:
      issue_transitions:
        enabled: false
      issue_comments:
        enabled: true
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
