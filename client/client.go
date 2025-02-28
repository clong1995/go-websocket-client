package client

import (
	"fmt"
	"github.com/clong1995/go-config"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var c = client{}

var wsClose = make(chan bool)

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
	wsClose <- true
}

func listen() {
	c.mu.RLock()
	defer func() {
		c.mu.RUnlock()
		_ = c.conn.Close()
	}()
	for {
		//当前客户端不需要接收消息
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func Connect() {
	go func() {
		for {
			select {
			case <-wsClose:
				fmt.Println("[websocket] closed")
				return
			default:
				err := connect()
				if err != nil {
					fmt.Println("[websocket] reconnecting")
					time.Sleep(1 * time.Second)
					continue
				}
				listen()
			}
		}
	}()

}
