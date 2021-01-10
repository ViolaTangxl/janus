package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Redis   RedisConfig   `yaml:"redis"`
	Session SessionConfig `yaml:"session"`
	ProxyCfg           ProxyConfig        `yaml:"proxy_config"`
	ProxyEntries       []ProxyEntry       `yaml:"-"` // 由 proxyCfg.files 指定的子配置文件填充
}

// ServerConfig config for server
type ServerConfig struct {
	Port string  `yaml:"port"`
	Mode RunMode `yaml:"mode"`
}

// RedisConfig config for redis service
type RedisConfig struct {
	Addrs      []string `yaml:"addrs"`
	MasterName string   `yaml:"master_name"`
	Failover   bool     `yaml:"failover"`
	Password   string   `yaml:"password"`
	DB         int      `yaml:"db"`
	Size       int      `yaml:"size"`
	Networt    string   `yaml:"network"`
	KeyPairs   string   `yaml:"key_pairs"`
}

type SessionConfig struct {
	MaxAge int `yaml:"max_age"`
}

func ParseConfig(file string) (*Config, error) {
	confData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(confData, conf)
	if err != nil {
		return nil, err
	}

	if conf.Server.Mode == "" {
		conf.Server.Mode = ProdMode
	}

	return conf, nil
}
