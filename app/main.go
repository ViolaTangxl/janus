package main

import (
	"flag"

	"github.com/ViolaTangxl/l-bridge/config"
	"github.com/ViolaTangxl/l-bridge/env"
	"github.com/sirupsen/logrus"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "conf", "config.yml", "config file path")
	flag.Parse()

	conf, err := config.ParseConfig(configPath)
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to load config")
		return
	}

	app, err := env.InitAppEngine(logrus.StandardLogger(), conf)
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to init engine")
		return
	}

	err = app.Run(conf.Server.Port)
	if err != nil {
		logrus.Errorf("l-bridge service run with err, err: %s", err)
		return
	}
}
