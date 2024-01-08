package main

import (
	"context"
	"fmt"
	"money/conf"
	"money/engine"
	"money/pkg/log"
	"money/pkg/mongodb"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
	err         error
)

func main() {
	f := kingpin.New(filepath.Base(os.Args[0]), "disk").Author("base cloud platform")

	f.Flag("config", "config.yaml").Default("").StringVar(&conf.ConfigPath)

	_, err := f.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse command failed, try command --help")
		f.Usage(os.Args[1:]) // auto exit(1)
	}

	if conf.ApolloPath == "" && conf.ConfigPath == "" {
		f.Usage(nil)
	}

	//创建一个服务
	conf.Init()
	cfg := conf.GetConfig()
	log.Init(&cfg.Log)
	_, err = mongodb.Init(ctx, &cfg.Mongodb)
	if err != nil {
		log.Fatalf("mongodb init failed, err: %s", err.Error())
	}
	engine.StartServer(ctx)

}
