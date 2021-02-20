package controllers

import (
	"github.com/ViolaTangxl/janus/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	config *config.Config
	logger logrus.FieldLogger
}

func NewProxyHandler(cfg *config.Config, logger logrus.FieldLogger) *Proxy {
	return &Proxy{
		config: cfg,
		logger: logger,
	}
}

// HandleProxyRequest 执行 Proxy 的主要逻辑
func (s *Proxy) HandleProxyRequest(ctx *gin.Context) {
	// 增加、修改请求
	// ReverseProxy
}
