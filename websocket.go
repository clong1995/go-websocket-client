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

// var ir = &isRun{}
var ws string

func init() {
	ws = config.Value("WEBSOCKET")
	run()
}

func run() {
	//连接
	conn := connect()
	//监听消息
	go listen(conn)
}

func listen(conn *websocket.Conn) {
	for {
		select {
		case <-wsClose:
			_ = conn.Close()
			log.Println("websocket closed")
			return
		default:
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				var syntaxError *json.SyntaxError
				if errors.As(err, &syntaxError) { //JSON解析错误
					continue
				}
				_ = conn.Close()
				log.Println("reconnecting...")
				time.Sleep(time.Second)
				run()
				return
			}

			//TODO 处理收到的消息
		}
	}
}

func connect() (conn *websocket.Conn) {
	var err error
	for {
		headers := http.Header{}
		headers.Add("a", "xxx") //这里时间要生新生成
		if conn, _, err = websocket.DefaultDialer.Dial(ws, headers); err != nil {
			log.Println(err)
			log.Println("Connection lost, reconnecting")
			time.Sleep(1 * time.Second)
		} else {
			log.Println("Connection connected")
			break
		}
	}
	return
}

func Close() {
	wsClose <- true
}
