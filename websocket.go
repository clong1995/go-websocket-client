package websocket

import (
	"encoding/json"
	"errors"
	"github.com/clong1995/go-config"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var wsClose = make(chan bool)

func init() {
	go run()
}

func run() {
	if <-wsClose {
		return
	}
	//连接
	conn := connect()
	//监听消息
	listen(conn)
}

func listen(conn *websocket.Conn) {
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			var syntaxError *json.SyntaxError
			if errors.As(err, &syntaxError) { //JSON解析错误
				continue
			}
			_ = conn.Close()
			go run()
			return
		}

		//TODO 处理收到的消息
	}
}

func connect() (conn *websocket.Conn) {
	var err error
	ws := config.Value("WEBSOCKET")
	for {
		headers := http.Header{}
		headers.Add("a", "xxx") //这里时间要冲新生成
		if conn, _, err = websocket.DefaultDialer.Dial(ws, headers); err == nil {
			log.Println(err)
			log.Println("Connection lost, reconnecting")
		} else {
			log.Println("Connection connected")
			break
		}
		time.Sleep(1 * time.Second)
	}
	return
}

func Close() {
	wsClose <- true
}
