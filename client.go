package wsrouter

import (
	"golang.org/x/net/websocket"
	"github.com/drdreyworld/sequence"
	"github.com/drdreyworld/events"
)

var clientId = sequence.Int{}

const outgoingQueueSize = 10000

func NewClient(ws *websocket.Conn) (cli Client) {
	cli.init(ws)
	return
}

type Client struct {
	clid int
	conn *websocket.Conn
	outc chan events.Event
}

func (c *Client) init(ws *websocket.Conn) {
	c.clid = clientId.GetNext()
	c.conn = ws
	c.outc = make(chan events.Event, outgoingQueueSize)
}

func (c *Client) SubscriberID() int {
	return c.clid
}

func (c *Client) Notify(event events.Event) {
	c.outc <- event
}

func (c *Client) ListenOutgoing() {
	for {
		select {
		case event := <-c.outc:
			websocket.JSON.Send(c.conn, Event{
				Code: event.GetID(),
				Data: event.GetData(),
			})
		}
	}
}

func (c *Client) Receive() (msg Event, err error) {
	err = websocket.JSON.Receive(c.conn, &msg)
	return
}