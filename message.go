package websocket

// 消息队列
var messageQueue = make(chan Message, 100)

// Message 只传递动作信息，不得传递其他数据
type Message struct {
	Subject string  `json:"s"`
	From    int64   //消息的来源用户,这个值不是发送方设置的，是注册成功后主动设置的。
	Target  []int64 `json:"t"` //消息的目的用户
	Payload string  `json:"p"`
}
