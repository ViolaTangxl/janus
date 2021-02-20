package config

import (
	"github.com/sirupsen/logrus"
)

var Global GlobalEnv

type GlobalEnv struct {
	Cfg           *Config
	ProxyCfg      *ProxyConfig
	Logger        logrus.FieldLogger
}
