package api

import (
	"strconv"
	"os"
	"strings"
)

// Config represents the configuration for a
// specific domain. Each domain can have multiple
// hostnames, including the root domain, where
// hostname is an empty string.
//
// The interval is the polling time (in seconds) for
// daemon mode.
type Config struct {
	Dev       bool     `json:"dev"`
	Domain    string   `json:"domain"`
	Hostnames []string `json:"hostnames"`
	Interval  int      `json:"interval"`
	Token     string   `json:"token"`
	Username  string   `json:"username"`
}

// LoadConfigs loads configurations from a file. The configuration
// is stored as an array of JSON serialized Config structs.
func LoadConfig() (Config) {
	devValue, _ := strconv.ParseBool(os.Getenv("NAME_DEV_MODE"))
	hostnameVals := strings.Split(os.Getenv("NAME_HOSTNAMES"), ",")
	domainValue := os.Getenv("NAME_DOMAIN")
	intervalValue, _ := strconv.ParseInt(os.Getenv("NAME_INTERVAL"), 10,32)
	tokenValue := os.Getenv("NAME_TOKEN")
	usernameValue := os.Getenv("NAME_USER")

	config := Config{
		Dev:       devValue,
		Domain:    domainValue,
		Hostnames: hostnameVals,
		Interval:  int(intervalValue),
		Token:     tokenValue,
		Username:  usernameValue,
	}

	return config
}
