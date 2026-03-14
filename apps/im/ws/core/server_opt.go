package core

import (
	"bufio"
	"math"
	"os"
	"strings"
	"time"

	"dicetales.com/pkg/sensitive"
	"github.com/zeromicro/go-zero/core/limit"
)

const (
	defaultMaxConnectionIdle   = time.Duration(math.MaxInt64)
	defaultGroupMsgConcurrency = 100
)

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication

	patten string

	// 消息ACK相关
	ack        AckType
	ackTimeout time.Duration
	// 最长连接空闲时间，超过这个时间没有任何消息（任意消息，包括心跳包），则关闭连接
	maxConnectionIdle time.Duration
	// 群聊最大消息转发数
	groupMsgConcurrency int
	// 敏感词过滤器，基于布隆过滤器实现
	sensitiveFilter *sensitive.SensitiveFilter
	// 消息限流器，针对每个用户的消息发送频率进行限制
	MsgLimiter *limit.TokenLimiter
}

func newServerOptions(opts ...ServerOptions) serverOption {
	o := serverOption{
		Authentication:      new(authentication),
		maxConnectionIdle:   defaultMaxConnectionIdle,
		patten:              "/ws",
		groupMsgConcurrency: defaultGroupMsgConcurrency,
	}

	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}

func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOption) {
		opt.patten = patten
	}
}

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *serverOption) {
		opt.ack = ack
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}

func WithServerSensitive() ServerOptions {
	return func(opt *serverOption) {
		if opt.sensitiveFilter == nil {
			// 静态加载敏感词
			words, err := LoadSensitiveWordsFromFile("apps/im/ws/etc/dev/sensitive.txt")
			if err != nil {
				panic("failed to load sensitive words: " + err.Error())
			}
			opt.sensitiveFilter = sensitive.NewSensitiveFilter(words)
		}
	}
}

func WithServerMsgLimiter() ServerOptions {
	return func(opt *serverOption) {
		if opt.MsgLimiter == nil {
			// 每个用户每秒最多发送5条消息
			opt.MsgLimiter = limit.NewTokenLimiter(5, 20, nil, "")
		}
	}
}

// 按行读取敏感词
func LoadSensitiveWordsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}
