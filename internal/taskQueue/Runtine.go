package taskQueue

import (
	"context"

	"github.com/ManInM00N/go-tool/goruntine"
	"github.com/devchat-ai/gopool"
)

var (
	SinglePool   gopool.GoPool
	TaskPool     *goruntine.TaskPool
	P            *goruntine.TaskPool
	Ctx          context.Context
	Cancel       context.CancelFunc
	IsClosed     = false
	WaitingTasks = int64(0)
)
