package client

import (
	"github.com/clong1995/go-config"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

var c = client{}

var wsClose = make(chan bool)

type client struct {
	conn *websocket.Conn
	mu   sync.RWMutex
}

func connect() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	ws := config.Value("WEBSOCKET")
	for {
		conn, _, err := websocket.DefaultDialer.Dial(ws, nil)
		if err != nil {
			log.Println("[websocket] lost, reconnecting: ", err)
			time.Sleep(1 * time.Second)
			continue
		}
		c.conn = conn
		log.Println("[websocket] connected: ", ws)
		break
	}
}

func Close() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_ = c.conn.Close()
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
			log.Println(err)
			log.Println("[websocket] reconnecting...")
			return
		}
	}
}

func Connect() {
	go func() {
		for {
			select {
			case <-wsClose:
				log.Println("[websocket] closed")
				return
			default:
				connect()
				listen()
			}
		}
	}()

}
