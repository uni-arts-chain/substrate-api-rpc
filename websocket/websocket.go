package websocket

import (
	"fmt"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"net/http"
	"time"
)

var wsEndPoint = ""

type WsConn interface {
	Dial(urlStr string, reqHeader http.Header)
	IsConnected() bool
	Close()
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, message []byte, err error)
	WriteJSON(v interface{}) error
	ReadJSON(v interface{}) error
	MarkUnusable()
	CloseAndReconnect()
}

func Init() (*PoolConn, error) {
	var err error
	if wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{
				KeepAliveTimeout: 30 * time.Second,
				WriteTimeout:     30 *time.Second,
				ReadTimeout:      time.Second * 30,
				NonVerbose:       true,
				HandshakeTimeout: 30 * time.Second}
			SubscribeConn.Dial(wsEndPoint, nil)
			return SubscribeConn, err
		}
		if wsPool, err = NewChannelPool(1, 25, factory); err != nil {
			fmt.Println("NewChannelPool", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := wsPool.Get()
	return conn, err
}

func SetEndpoint(endpoint string) {
	wsEndPoint = endpoint
}

func Close() {
	if wsPool != nil {
		wsPool.Close()
	}
}
