package websocket

import (
	"math"
	"time"
)

const (
	defaultMaxConnectionIdle  = time.Duration(math.MaxInt64)
	defaultAckTimeout         = 3 * time.Second
	defaultGroupMsgConurrency = 100
)
