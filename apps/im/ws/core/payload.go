package core

type MType int

const (
	TextMType MType = iota
)

type ChatType int

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)

type (
	ContentWrapper struct {
		MType   `mapstructure:"mType"`
		ID      string `mapstructure:"id"`
		Content string `mapstructure:"content"`
	}

	ChatMessage struct {
		ConversationId string `mapstructure:"conversationId"`
		SendId         string `mapstructure:"sendId"`
		RecvId         string `mapstructure:"recvId"`
		SendTime       int64  `mapstructure:"sendTime"`

		ChatType `mapstructure:"chatType"`
		Wrapper  ContentWrapper `mapstructure:"content"`

		RecvIds []string `mapstructure:"recvIds"`
	}
)
