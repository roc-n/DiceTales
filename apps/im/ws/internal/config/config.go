package config

import "github.com/zeromicro/go-zero/core/service"

type Config struct {
	Snowflake struct {
		Node int64
	}

	service.ServiceConf

	ListenOn string

	JwtAuth struct {
		AccessSecret string
	}

	Mongo struct {
		Url string
		Db  string
	}

	MessageTransfer struct {
		Topic string
		Addrs []string
	}

	SensitiveWordsPath string
}
