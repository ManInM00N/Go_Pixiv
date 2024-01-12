package init

import (
	"context"
	"github.com/ManInM00N/go-tool/goruntine"
	"github.com/devchat-ai/gopool"
)

var (
	SinglePool gopool.GoPool
	P          gopool.GoPool
	TaskPool   *goruntine.GoPool
	Ctx        context.Context
	Cancel     context.CancelFunc
	IsClosed   = false
)
