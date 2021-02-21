package middlewares

import (
	"github.com/ViolaTangxl/janus/utils"
	"github.com/gin-gonic/gin"
)

const (
	ctxKeyReqID    = "reqid"
	headerKeyReqID = "X-Reqid"
)

func newReqid() string {
	return utils.GenRandomString(16)
}

// GetReqidMiddleware gets or gens reqid and sets to ctx
func GetReqidMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqid := ctx.GetHeader(headerKeyReqID)
		if reqid == "" {
			reqid = newReqid()
		}

		ctx.Set(ctxKeyReqID, reqid)
		ctx.Header(headerKeyReqID, reqid)
		ctx.Next()
	}
}

// GetReqid gets reqid from ctx
func GetReqid(ctx *gin.Context) string {
	return ctx.MustGet(ctxKeyReqID).(string)
}
