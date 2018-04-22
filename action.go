package wsrouter

import "context"

type ActionFunc func(ctx context.Context) (interface{}, error)

type Action struct {
	action  ActionFunc
	context context.Context
}

type Actions map[string]Action
