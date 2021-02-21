package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/ViolaTangxl/janus/config"
	"github.com/ViolaTangxl/janus/controllers/middlewares"
	"github.com/gin-contrib/sessions"
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
	log := middlewares.GetLogger(ctx)
	pathArray, name, err := s.vaildRequestPath(ctx)
	if err != nil {
		logrus.Errorf("<proxy.HandleProxyRequest> vaildRequestPath() failed, RespErr:%d, err:%s.", http.StatusNotFound, err)
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	// 增加、修改请求
	targetInfo, host, err := s.getTargetAndHost(ctx, name, pathArray)
	if err != nil {
		log.Errorf("<proxy.HandleProxyRequest> getTargetAndHost() failed, RespErr:%d, err:%s.", http.StatusNotFound, err)
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	err = s.addParam(ctx, targetInfo)
	if err != nil {
		log.Errorf("<proxy.HandleProxyRequest> addParam(%v) failed, RespErr:%d, err:%s.", targetInfo, http.StatusNotFound, err)
		ctx.JSON(http.StatusNotFound, nil)
		return
	}

	targetURL, err := url.Parse(host)
	if err != nil {
		log.Errorf("<proxy.HandleProxyRequest> url.Parse(%s) failed, RespErr:%d, err:%s.", host, http.StatusNotFound, err)
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	targetQuery := targetURL.RawQuery

	// ReverseProxy
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = targetInfo.Path
			req.Host = targetURL.Host
			req.Header = ctx.Request.Header
			req.Header.Set("Host", targetURL.Host)
			req.Header.Del("Accept-Encoding")
			req.Header.Set("X-Reqid", middlewares.GetReqid(ctx))
			req.Method = string(targetInfo.Method)
			req.Body = ctx.Request.Body
			req.Form = ctx.Request.Form
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
		},
	}
	s.logger.Infof("<proxy.HandleProxyRequest> request path: %s proxy to ==> : %s%s, "+
		"request Method: %s, rawQuery: %s, reqid: %s, refer: %s.",
		ctx.Request.URL.Path,
		host,
		targetInfo.Path,
		ctx.Request.Method,
		ctx.Request.URL.RawQuery,
		middlewares.GetReqid(ctx),
		ctx.Request.Referer())
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

// vaildRequestPath 判断是否是有效请求(必须以/api/proxy 开头)
func (s *Proxy) vaildRequestPath(ctx *gin.Context) ([]string, string, error) {
	var (
		path    = ctx.Request.URL.Path
		matches = strings.Split(path, "/")
	)
	// 只有 /api/proxy
	if len(matches) <= 3 {
		err := errors.New("invalid path")
		return nil, "", err
	}

	name := matches[3]

	// path 只有/api/proxy/
	if len(name) == 0 {
		err := errors.New("invalid path")
		return nil, "", err
	}
	return matches[4:], name, nil
}

// getTargetAndHost 自定义请求参数
func (s *Proxy) getTargetAndHost(ctx *gin.Context, name string, pathArray []string) (*config.Match, string, error) {
	var (
		host      string
		matchInfo config.Match
		err       error
	)
	for _, proxyEntry := range s.config.ProxyEntries {
		if strings.TrimSpace(proxyEntry.Name) != name {
			continue
		}

		host = proxyEntry.Target
		for _, match := range proxyEntry.Matches {
			if isSameMethod(ctx.Request.Method, string(match.Method)) &&
				isSamePath(pathArray, strings.TrimSpace(match.Path)) {
				matchInfo.Path = match.Path
				matchInfo.Method = ctx.Request.Method
				matchInfo.Params = match.Params
				return &matchInfo, host, nil
			}
		}
	}

	err = errors.New("not find")
	return &matchInfo, host, err
}

func isSameMethod(targetMethod, matchInfoMethod string) bool {
	if matchInfoMethod == "*" {
		return true
	}
	if strings.ToUpper(targetMethod) == strings.ToUpper(matchInfoMethod) {
		return true
	}

	return false
}

func isSamePath(targetPathArray []string, matchInfo string) bool {
	// 支持配置/*或者*
	if matchInfo == "*" || matchInfo == "/*" {
		return true
	}
	if matchInfo == "/" || !strings.HasPrefix(matchInfo, "/") {
		return false
	}

	matchInfoPathArray := strings.Split(matchInfo, "/")[1:]

	matchLen := len(matchInfoPathArray)
	targetLen := len(targetPathArray)

	if (targetLen < matchLen && matchInfoPathArray[matchLen-1] != "*") ||
		(targetLen < matchLen-1 && matchInfoPathArray[matchLen-1] == "*") {
		return false
	}

	for i, targetPath := range targetPathArray {
		if i >= matchLen {
			return false
		}

		if strings.HasPrefix(matchInfoPathArray[i], ":") {
			continue
		}

		if matchInfoPathArray[i] == "*" {
			return true
		}
		if targetPath != matchInfoPathArray[i] {
			return false
		}
	}

	return true
}

func (s *Proxy) addParam(ctx *gin.Context, matchInfo *config.Match) (err error) {
	var (
		session    = sessions.Default(ctx)
		paramMap   = make(map[string]interface{})
		reqPath    = strings.Split(ctx.Request.URL.Path, "/")[4:]
		targetPath = fmt.Sprintf("/%s", strings.Join(reqPath, "/"))
		newBody    []byte
	)

	if ctx.Request.Body != nil {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			s.logger.Errorf("ioutil.ReadAll err: ", err)
			return err
		}
		if len(body) > 0 {
			err = json.Unmarshal(body, &paramMap)
			if err != nil {
				err = errors.New("unmarshal failed")
				return err
			}
		}
	}

	for _, paramInfo := range matchInfo.Params {
		paramKey, paramValue, err := s.getParamMessage(session, paramInfo)
		if err != nil {
			return err
		}

		switch paramInfo.Location {
		case config.ParamLocationUrlParam:
			value := ctx.Request.URL.Query()
			value.Add(paramKey, fmt.Sprintf("%v", paramValue))
			ctx.Request.URL.RawQuery = value.Encode()
		case config.ParamLocationBody:
			paramMap[paramKey] = paramValue
		case config.ParamLocationHeader:
			ctx.Request.Header.Set(paramKey, fmt.Sprintf("%v", paramValue))
		case config.ParamLocationUrlPath:
			targetPath, err = updateHost(reqPath, matchInfo, paramKey, fmt.Sprintf("%v", paramValue))
			if err != nil {
				return err
			}
		}
	}

	if len(paramMap) > 0 {
		newBody, _ = json.Marshal(paramMap)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))
		ctx.Request.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
		ctx.Request.ContentLength = int64(len(newBody))
	}

	// 修改 targetPath
	matchInfo.Path = targetPath
	return nil
}

// updateHost 针对 url_path 的情况修改实际请求中的参数
func updateHost(reqPaths []string, match *config.Match, paramKey, paramValue string) (string, error) {
	// NOTE：假设 match.Paths = /cps/:uid,则 split 后为["","cps",":uid"]
	for i, path := range strings.Split(match.Path, "/")[1:] {
		splitArray := strings.Split(path, ":")
		if len(splitArray) == 2 && splitArray[0] == "" && paramKey == splitArray[1] {
			// 虽然不会出现请求的 url 比配置的 path 短的情况，但是防止 panic 还是检查一下
			if len(reqPaths) <= i {
				return "", errors.New("invaild request")
			}
			reqPaths[i] = paramValue
		}
	}

	return fmt.Sprintf("/%s", strings.Join(reqPaths, "/")), nil
}

func (s *Proxy) getParamMessage(session sessions.Session, param config.Param) (paramKey string, paramValue interface{}, err error) {
	if param.SessionKey != "" {
		if param.Rename != "" {
			paramKey = param.Rename
		} else {
			paramKey = param.SessionKey
		}

		sessionValue := session.Get(param.SessionKey)
		// 仅当 session 和 custom value均没有值时,才报错返回
		if sessionValue == nil && param.CustomValue == nil {
			err = errors.New("session get value failed and custom value is nil")
			return
		}

		// sessionValue 优先级高于custom value,优先从 session中填充value
		if sessionValue != nil {
			paramValue = sessionValue
		} else {
			paramValue = param.CustomValue
		}
	} else {
		paramKey = param.Rename
		if param.CustomValue == nil {
			s.logger.Errorf("<Proxy.getParamMessage> paramKey: %s, customValue is nil", paramKey)
			err = errors.New("customValue is nil")
			return
		}
		paramValue = param.CustomValue
	}

	return
}
