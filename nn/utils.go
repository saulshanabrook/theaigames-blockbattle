package nn

import "github.com/gorilla/websocket"

func newConn(host, path string) *websocket.Conn {
	d := websocket.Dialer{}
	u := "ws://" + host + path
	conn, _, err := d.Dial(u, nil)
	handle(err)
	return conn
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
