package websocket

import "errors"

func Send(msg Message) (err error) {
	/*if err = conn.WriteJSON(msg); err != nil {
		log.Println(err)
		return
	}*/
	select {
	case messageQueue <- msg:
	default:
		err = errors.New("message queue is full")
	}
	return
}
