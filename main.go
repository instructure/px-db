package main

import (
	"os"
	"strings"

	"github.com/instructure/px-db/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("PX_DB_LOG_LEVEL", "INFO")

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	log.SetReportCaller(true)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log debug or info levels
	logLevel := viper.GetString("px_db_log_level")
	log.Println("Using Log Level:", logLevel)
	if strings.ToLower(logLevel) == "info" {
		log.SetLevel(log.InfoLevel)
	} else if strings.ToLower(logLevel) == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.Fatal("Invalid log level specified - valid values: 'INFO' and 'DEBUG'")
	}
}

func main() {
	rootCommand := cmd.NewCmdRoot()
	if err := rootCommand.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
