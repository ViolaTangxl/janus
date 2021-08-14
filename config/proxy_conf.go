package config

import (
	"errors"
)

type ParamLocation string

const (
	ParamLocationUrlParam ParamLocation = "url_param"
	ParamLocationBody     ParamLocation = "body"
	ParamLocationHeader   ParamLocation = "header"
	ParamLocationUrlPath  ParamLocation = "url_path"
)

var (
	ParseProxyConfErr              = errors.New("parse proxyConfig failed.")
	ProxyNameIsNil                 = errors.New("proxyConfig name is nil.")
	ProxyTargetIsNil               = errors.New("proxyConfig target is nil.")
	ProxyMatchesPathIsNil          = errors.New("proxyConfig matches path is nil.")
	ProxyMatchesMethodIsNil        = errors.New("proxyConfig matches method is nil.")
	ProxyMatchesParamLacationIsNil = errors.New("proxyConfig matches param location is nil.")
)
