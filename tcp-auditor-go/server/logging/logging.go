package logging

import (
	"bitbucket.tylertech.com/spy/scm/tcp-auditor/server/config"
	log "github.com/sirupsen/logrus"
)

// Logger is an instance of logrus logger
var Logger = log.New()

// InitLogger sets up logging level, and log formatting
func InitLogger() {
	c := config.GetConfig()
	ll, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Error("Unable to determine log level from configuration; defaulting to debug log level")
		ll = log.DebugLevel
	}
	log.SetLevel(ll)
	log.SetFormatter(&log.JSONFormatter{})
}
