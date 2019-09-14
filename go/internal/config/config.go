package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	DualPass             = "settings.dual_pass.enabled"
	DualPassWaitDuration = "settings.dual_pass.wait_duration"
	IssueComments        = "settings.issues.enable_comment"
	IssueTransitions     = "settings.issues.enable_transition"
	IssueConflictStatus  = "settings.issues.conflict_status"
	Jira                 = "settings.jira.enabled"
	JiraHost             = "settings.jira.host"
	JiraProjectName      = "settings.jira.project_name"
	JiraUser             = "settings.jira.user"
)

func Reset() {
	viper.Reset()
}

// CheckMessage is a helper function for building error messages related to configuration settings
func CheckMessage(details ...string) (msg string) {

	if len(details) == 0 {
		return
	}

	settingName := details[0]
	if len(details) == 1 {
		msg = fmt.Sprintf("check config.yaml: '%s'", settingName)
	}

	if len(details) > 1 {
		msg = fmt.Sprintf("check config.yaml: '%s'. %s", settingName, details[1])
	}

	return
}

func Initialize() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./github-actions/prwatch-action/")
	viper.AddConfigPath("../github-actions/prwatch-action/")
	viper.AddConfigPath("../../github-actions/prwatch-action/")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read configuration: %s", err)
	}

	viper.SetDefault(DualPass, true)
	viper.SetDefault(DualPassWaitDuration, "60s")
	viper.SetDefault(IssueComments, true)
	viper.SetDefault(IssueTransitions, true)
	viper.SetDefault(Jira, true)
}

func GlobalDisable(setting string) {

	viper.Set(setting, false)
}

func GlobalEnable(setting string) {

	viper.Set(setting, true)
}

func GlobalSet(setting, value string) {

	viper.Set(setting, value)
}

func SetEnv(envVar, value string) {

	viper.Set(envVar, value)
}

func GetEnv(envVar string) string {

	return viper.GetString(envVar)
}

func UserDisableSetting(user, setting string) {

	viper.Set(userSettingName(user, setting), false)
}

func UserEnable(user, setting string) {

	viper.Set(userSettingName(user, setting), true)
}

func GetString(setting string) string {

	return viper.GetString(setting)
}

func GetBool(setting string) bool {

	return viper.GetBool(setting)
}

func GetDuration(setting string) time.Duration {

	return viper.GetDuration(setting)
}

func IsSet(setting string) bool {

	return viper.IsSet(setting)
}

func SettingEnabled(setting string) bool {

	return viper.GetBool(setting)
}

func UserSettingEnabled(user, setting string) bool {

	// check if it's configured globally first, superceding user settings
	if IsSet(setting) {
		return SettingEnabled(setting)
	}

	userSetting := userSettingName(user, setting)

	return viper.GetBool(userSetting)
}

func userSettingName(user, setting string) string {
	return fmt.Sprintf("users.%s.%s", user, setting)
}
