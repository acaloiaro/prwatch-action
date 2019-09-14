package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

func Reset() {
	viper.Reset()
}

func Initialize() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./github-actions/prwatch-action/")
	viper.AddConfigPath("../github-actions/prwatch-action/")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read configuration: %s", err)
	}

}

func GlobalDisable(group string) {
	viper.Set(globalSettingName(group, "enabled"), false)
}

func GlobalEnable(settingGroup string) {
	viper.Set(globalSettingName(settingGroup, "enabled"), true)
}

func GlobalSet(group, setting, value string) {
	viper.Set(globalSettingName(group, setting), value)
}

func SetEnv(envVar, value string) {
	viper.Set(envVar, value)
}

func GetEnv(envVar string) string {
	log.Println("Fetching env var", envVar)
	val := viper.GetString(envVar)
	log.Println("got value:", val)
	return val
}

func UserDisable(user, setting string) {
	viper.Set(userSettingName(user, setting, "enabled"), false)
}

func UserEnable(user, setting string) {
	viper.Set(userSettingName(user, setting, "enabled"), true)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetString(group, setting string) string {
	if SettingEnabled(group) {
		return viper.GetString(globalSettingName(group, setting))
	}

	return ""
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetDuration(group, setting string) time.Duration {
	return viper.GetDuration(globalSettingName(group, setting))
}

func IsSetGlobal(group, setting string) bool {
	return viper.IsSet(globalSettingName(group, setting))
}

func SettingEnabled(query ...string) bool {
	if len(query) == 1 {

		setting := query[0]

		return viper.GetBool(globalSettingName(setting, "enabled"))

	} else if len(query) == 2 {

		user := query[0]
		setting := query[1]

		// check if it's disabled globally first, superceding user settings
		if !SettingEnabled(setting) {
			return false
		}

		full := userSettingName(user, setting, "enabled")
		log.Printf("User: %s, setting: %s, fqsn: %s", user, setting, full)
		return viper.GetBool(full)
	}

	return false
}

func globalSettingName(group, setting string) string {
	return fmt.Sprintf("settings.%s.%s", group, setting)
}

func userSettingName(user, group, setting string) string {
	return fmt.Sprintf("users.%s.settings.%s.%s", user, group, setting)
}

func parseDuration(durationString, configVarName string) (duration time.Duration) {
	var err error

	duration, err = time.ParseDuration(durationString)
	if err != nil {
		duration = time.Duration(0 * time.Second)
	}

	return
}
