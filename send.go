package websocket

func Send(msg Message) (err error) {
	messageQueue <- msg
	/*if err = conn.WriteJSON(msg); err != nil {
		log.Println(err)
		return
	}*/
	return
}
