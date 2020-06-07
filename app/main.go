package main

import (
	"flag"

	"github.com/ViolaTangxl/l-bridge/config"
	"github.com/sirupsen/logrus"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "conf", "config.yml", "config file path")
	flag.Parse()

	_, err := config.ParseConfig(configPath)
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to load config")
		return
	}
}
