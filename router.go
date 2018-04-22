package wsrouter

import (
	"context"
	"errors"
	"github.com/drdreyworld/events"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
)

func CreateRouter(events *events.Events) *Router {
	return &Router{
		actions: Actions{},
		events:  events,
	}
}

type Router struct {
	actions Actions
	events  *events.Events
}

func (r *Router) BindAction(id string, action ActionFunc, ctx context.Context) {
	r.actions[id] = Action{
		action:  action,
		context: ctx,
	}
}

func (r *Router) Execute(msg Event, cli *Client) (res interface{}, err error) {
	a, ok := r.actions[msg.Code]
	if ok {
		ctx := context.WithValue(a.context, "client", cli)
		ctx = context.WithValue(ctx, "params", msg.Data)

		res, err = a.action(ctx)
	} else {
		err = errors.New("Route not matched: " + msg.Code)
	}

	return
}

func (r *Router) onClientConnect(ws *websocket.Conn) {
	cli := NewClient(ws)

	defer r.onClientDisconnect(cli)
	go cli.ListenOutgoing()

	for {
		msg, err := cli.Receive()
		if err == io.EOF {
			return
		} else if err != nil {
			cli.Notify(CreateErrorEvent(msg.Code, err.Error()))
		} else {
			go func() {
				res, err := r.Execute(msg, &cli)
				if err != nil {
					cli.Notify(CreateErrorEvent(msg.Code, err.Error()))
				} else if res != nil {
					cli.Notify(CreateEvent(msg.Code, res))
				}
			}()
		}
	}
}

func (r *Router) onClientDisconnect(cli Client) {
	if err := cli.conn.Close(); err != nil {
		log.Println(err)
	}
}

func (r *Router) Bind(addr string) {
	http.Handle(addr, websocket.Handler(r.onClientConnect))

	r.BindAction(
		"subscribe",
		func(ctx context.Context) (interface{}, error) {
			client := ctx.Value("client").(events.Subscriber)
			params := ctx.Value("params").(map[string]interface{})
			if client != nil {
				if params != nil {
					r.events.Subscribe(params["event"].(string), client)
					return params, nil
				}
				return nil, errors.New("can't get params from context")
			}
			return nil, errors.New("can't get client from context")
		},
		context.Background(),
	)

	r.BindAction(
		"unsubscribe",
		func(ctx context.Context) (interface{}, error) {
			client := ctx.Value("client").(events.Subscriber)
			params := ctx.Value("params").(map[string]interface{})
			if client != nil {
				if params != nil {
					r.events.Unsubscribe(params["event"].(string), client)
					return params, nil
				}
				return nil, errors.New("can't get params from context")
			}
			return nil, errors.New("can't get client from context")
		},
		context.Background(),
	)
}
