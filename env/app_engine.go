package env

import (
	"github.com/ViolaTangxl/janus/config"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitAppEngine(l *logrus.Logger, cfg *config.Config) (*gin.Engine, error) {
	gin.SetMode(string(cfg.Server.Mode))

	app := gin.New()
	store, err := redis.NewStore(
		cfg.Redis.Size,
		cfg.Redis.Networt,
		cfg.Redis.Addrs[0],
		cfg.Redis.Password,
		[]byte(cfg.Redis.KeyPairs),
	)

	if err != nil {
		l.Errorf("<appengine.InitAppEngine> redis.NewStore() failed, err: %s.", err)
		return nil, err
	}
	store.Options(sessions.Options{MaxAge: cfg.Session.MaxAge})

	app.Use(
		gin.Recovery(),
	)

	return app, nil
}
