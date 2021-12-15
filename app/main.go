package main

import (
	"errors"
	"flag"
	"path/filepath"

	"github.com/ViolaTangxl/janus/config"
	"github.com/ViolaTangxl/janus/env"

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
	err = buildProxyEntries(conf)
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to build proxy entries")
		return
	}

	app, err := env.InitAppEngine(logrus.StandardLogger(), conf)
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to init engine")
		return
	}

	env.InitRouters(app)
	err = app.Run(conf.Server.Port)
	if err != nil {
		logrus.Errorf("l-bridge service run with err, err: %s", err)
		return
	}
}

// buildProxyEntries 加载多个配置文件
func buildProxyEntries(conf *config.Config) error {
	absPath, err := filepath.Abs(conf.ProxyCfg.Files)
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to read dir")
		return err
	}

	matches, err := filepath.Glob(absPath)
	if err != nil {
		logrus.WithField("err", err).Fatal("failed to match proxy cfg")
		return err
	}
	if len(matches) <= 0 {
		err = errors.New("proxy cfg not found")
		logrus.WithField("err", err).Fatal("failed to load proxy cfg")
		return err
	}

	proxyEntries := make([]config.ProxyEntry, 0)

	for _, match := range matches {
		pxy, err1 := config.ParseProxyEntry(match)
		if err1 != nil {
			logrus.WithField("err", err1).Fatal("failed to load proxy config")
			return err
		}
		proxyEntries = append(proxyEntries, pxy.ProxyEntries...)
	}
	conf.ProxyEntries = proxyEntries

	return nil
}
