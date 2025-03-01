package client

import (
	"context"
	"fmt"
	"github.com/clong1995/go-config"
	"github.com/clong1995/go-websocket-client/message"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

var c = client{}

var serverClose = make(chan bool)

var ctx, cancel = context.WithCancel(context.Background())

type client struct {
	conn *websocket.Conn
	mu   sync.RWMutex
}

func connect() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	ws := config.Value("WEBSOCKET")
	c.conn, _, err = websocket.DefaultDialer.Dial(ws, nil)
	if err != nil {
		return
	}
	fmt.Println("[websocket] connected: ", ws)
	return
}

func Close() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	cancel()
	close(queue)
	serverClose <- true
}

func listen() {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	defer func() {
		_ = conn.Close()
	}()
	for {
		//当前客户端不需要接收消息
		if _, _, err := conn.ReadMessage(); err != nil {
			//log.Println(err)
			return
		}
	}
}

func Connect() {
	go func() {
		for {
			select {
			case <-serverClose:
				fmt.Println("[websocket] exited!")
				return
			default:
				err := connect()
				if err == nil {
					listen()
				}
				//<=== 延时1秒后重试，在这里为防止`listen()`关闭后，立即进入 `connect()`，
				// 造成`serverClose <- true`没有及时写入而无法关闭
				time.Sleep(1 * time.Second)
				fmt.Println("[websocket] reconnecting")
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-queue:
				if !ok {
					return
				}
				sem <- struct{}{}
				go send(msg)
			}
		}
	}()
}

var sem = make(chan struct{}, 10)

// 发送给WS服务端
func send(msg message.Msg) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	defer func() { <-sem }()

	err := conn.WriteJSON(msg)
	if err != nil {
		log.Println(err)
		select {
		case queue <- msg: //发生问题后尝试在写回消息队列
		default:
		}
		return
	}
}
