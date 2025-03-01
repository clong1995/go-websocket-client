package client

import (
	"errors"
	"github.com/clong1995/go-websocket-client/message"
	"log"
)

// Queue 消息队列
var queue = make(chan message.Msg, 1000)

func Send(msg message.Msg) (err error) {
	select {
	case queue <- msg:
	default:
		err = errors.New("message queue is full")
		log.Println(err)
	}
	return
}
