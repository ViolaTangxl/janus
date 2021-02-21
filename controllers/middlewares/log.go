package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const ctxKeyReqLogger = "logger"

// GetLoggerMiddleware 向 ctx 中注入 logger
// 依赖 reqid middleware，需要先 use reqid middleware，再 use logger middleware
func GetLoggerMiddleware(l *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqid := GetReqid(ctx)
		logger := l.WithFields(logrus.Fields{
			"reqid": reqid,
		})
		ctx.Set(ctxKeyReqLogger, logger)
		ctx.Next()
	}
}

// GetLogger 用于获取 ctx 中的 logger
func GetLogger(ctx *gin.Context) *logrus.Entry {
	return ctx.MustGet(ctxKeyReqLogger).(*logrus.Entry)
}

// GetLogReqMiddleware 会记录每一条 request 的基本信息
// 依赖 logger middleware，需要先 use logger middleware，再 use log req middleware
func GetLogReqMiddleware(l *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			startTime = time.Now()
			latency   time.Duration
			status    int
			logger    = GetLogger(ctx)
		)

		clientIP := ctx.ClientIP()

		info := fmt.Sprintf(
			"request: %s %s %s %s %s %s",
			clientIP, ctx.Request.Method,
			ctx.Request.URL.Path, ctx.Request.Header.Get("Origin"), ctx.Request.Header.Get("Referer"),
			ctx.Request.Header.Get("X-Original-Forwarded-For"), // for k8s
		)
		logger.Info(info)

		// Process request
		ctx.Next()

		latency = time.Now().Sub(startTime)
		status = ctx.Writer.Status()

		info = fmt.Sprintf(
			"response: %s %s %s %s %s %s %d %s",
			clientIP, ctx.Request.Method,
			ctx.Request.URL.Path, ctx.Request.Header.Get("Origin"), ctx.Request.Header.Get("Referer"),
			ctx.Request.Header.Get("X-Original-Forwarded-For"), // for k8s
			status, latency,
		)
		if status/100 == 2 {
			logger.Info(info)
		} else {
			logger.Error(info)
		}
	}
}
