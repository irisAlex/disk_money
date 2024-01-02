package main

import (
	"context"
	"money/engine"
)

var ctx, cancel = context.WithCancel(context.Background())

func main() {
	//创建一个服务
	engine.StartServer(ctx)
}
