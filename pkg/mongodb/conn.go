package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbName    string
	ctx       context.Context
	MgoClient *mongo.Client // allow direct call

	replayCount       = 5
	replayTimeoutWait = 5 * time.Second
	timeoutWait       = 25 * time.Second
	connectTimeout    = 10 * time.Second

	errMongodbConnect = errors.New("mongodb connect error ")
)

// Scheme constants
const (
	SchemeMongoDB    = "mongodb"
	SchemeMongoDBSRV = "mongodb+srv"
)

type Config struct {
	Endpoint   string `toml:"endpoint"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	DB         string `toml:"db"`
	Source     string `toml:"source"`
	MaxPool    uint64 `toml:"max_pool"`
	ReplicaSet string `toml:"replica_set" yaml:"replica_set"`
}

func Init(hctx context.Context, cfg *Config) (func(), error) {
	var (
		warning = make(chan error)
		count   = new(int)
		serious = make(chan error)
		sigIn   = make(chan struct{})

		credential    options.Credential
		clientOptions *options.ClientOptions
	)

	dbName = cfg.DB
	ctx = hctx

	if !strings.HasPrefix(cfg.Endpoint, SchemeMongoDBSRV+"://") && !strings.HasPrefix(cfg.Endpoint, SchemeMongoDB+"://") {
		cfg.Endpoint = fmt.Sprintf("%s://%s", SchemeMongoDB, cfg.Endpoint)
	}
	m := &Mongodb{
		warning:  warning,
		serious:  serious,
		sigIn:    sigIn,
		count:    count,
		endpoint: cfg.Endpoint,
	}

	if cfg.User != "" && cfg.Password != "" {
		credential = options.Credential{
			AuthSource: cfg.Source,
			Username:   cfg.User,
			Password:   cfg.Password,
		}
		clientOptions = options.Client().ApplyURI(cfg.Endpoint).SetAuth(credential)
	} else {
		clientOptions = options.Client().ApplyURI(cfg.Endpoint)
	}

	if cfg.MaxPool != 0 {
		clientOptions.SetMaxPoolSize(cfg.MaxPool)
	}
	if cfg.ReplicaSet != "" {
		clientOptions.SetReplicaSet(cfg.ReplicaSet)
	}

	go m.HandleMongodbSession(clientOptions)

	select {
	case <-m.sigIn:
		close(m.serious)
		close(m.warning)
		close(m.sigIn)
		return Shutdown, nil

	case err := <-m.serious:
		return Shutdown, err

	case <-time.After(timeoutWait):
		return Shutdown, errMongodbConnect
	}
}

func Shutdown() {
	_ = MgoClient.Disconnect(ctx)
}

type Mongodb struct {
	warning  chan error
	serious  chan error
	sigIn    chan struct{}
	count    *int
	endpoint string
}

func (m *Mongodb) HandleMongodbSession(clientOptions *options.ClientOptions) {
	var (
		ctx = context.Background()
		err error
	)

	for {
		MgoClient, err = mongo.NewClient(
			clientOptions, clientOptions.SetConnectTimeout(connectTimeout),
		)
		if err != nil {
			log.Printf("mongo.NewClient error: %s ", err.Error())
			m.serious <- err
			return
		}

		if err := MgoClient.Connect(ctx); err != nil {
			log.Printf("mongodb connect error: %s", err.Error())
			m.serious <- err
			return
		}

		err = MgoClient.Ping(ctx, readpref.Primary())
		if err != nil && *m.count < replayCount {
			// 连接mongodb失败,重试5次
			*m.count += 1
			time.Sleep(replayTimeoutWait)
		} else if err != nil {
			m.serious <- err
			return
		} else {
			m.sigIn <- struct{}{}
			return
		}
	}
}

func GetMongoClient() *mongo.Database {
	return MgoClient.Database(dbName)
}

func GetMongoCollection(table string) *mongo.Collection {
	return MgoClient.Database(dbName).Collection(table)
}
