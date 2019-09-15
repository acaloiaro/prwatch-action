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

func GlobalUnset(setting string) {
	viper.Set(setting, nil)
}

func SetEnv(envVar, value string) {
	viper.Set(envVar, value)
}

func GetEnv(envVar string) string {
	return viper.GetString(envVar)
}

func UserDisable(user, setting string) {

	viper.Set(userSettingName(user, setting), false)
}

func UserEnable(user, setting string) {

	viper.Set(userSettingName(user, setting), true)
}

func GetString(setting string) string {

	return viper.GetString(setting)
}

func GetBool(setting string) bool {

	log.Println("Getting", setting)
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

	userSetting := userSettingName(user, setting)

	// check if it's configured globally first, superceding user settings
	if IsSet(setting) {
		return SettingEnabled(setting)
	}

	return viper.GetBool(userSetting)
}

func userSettingName(user, setting string) string {
	return fmt.Sprintf("users.%s.%s", user, setting)
}
