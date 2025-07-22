package taskQueue

import (
	"context"

	"github.com/ManInM00N/go-tool/goruntine"
	"github.com/devchat-ai/gopool"
)

var (
	SinglePool     gopool.GoPool
	P              gopool.GoPool
	RankloadPool   gopool.GoPool
	FollowLoadPool gopool.GoPool
	FollowPool     *goruntine.GoPool
	TaskPool       *goruntine.GoPool
	RankPool       *goruntine.GoPool
	Ctx            context.Context
	Cancel         context.CancelFunc
	IsClosed       = false
	WaitingTasks   = int64(0)
)
