package env

import (
	"github.com/ViolaTangxl/janus/config"
	"github.com/ViolaTangxl/janus/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouters(app *gin.Engine) {
	proxyHandler := controllers.NewProxyHandler(
		config.Global.Cfg,
		config.Global.Logger,
	)

	proxy := app.Group("/api/proxy")
	{
		proxy.Any("/*proxy", proxyHandler.HandleProxyRequest)
	}
}
