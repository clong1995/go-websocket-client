package websocket

import (
	"errors"
	"github.com/clong1995/go-websocket-client/message"
	"log"
)

func Send(msg message.Msg) (err error) {
	select {
	case message.Queue <- msg:
	default:
		err = errors.New("message queue is full")
		log.Println(err)
	}
	return
}
