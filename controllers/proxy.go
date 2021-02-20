package controllers

import (
	"github.com/sirupsen/logrus"
	"janus/config"
)

type Proxy struct {
	config *config.Config
	logger *logrus.FieldLogger
}

func NewProxyHandler(cfg *config.Config, logger *logrus.FieldLogger) *Proxy {
	return &Proxy{
		config:cfg,
		logger:logger,
	}
}
