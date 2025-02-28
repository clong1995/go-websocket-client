package websocket

import "github.com/clong1995/go-websocket-client/client"

func init() {
	client.Connect()
}

func Close() {
	client.Close()
}
