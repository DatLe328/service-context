package mongoc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	sctx "github.com/DatLe328/service-context"
	"github.com/DatLe328/service-context/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoOpt struct {
	uri                    string
	database               string
	maxPoolSize            uint64
	minPoolSize            uint64
	maxConnIdleTime        int
	connectTimeout         int
	serverSelectionTimeout int
}

type mongoDB struct {
	id       string
	prefix   string
	logger   logger.Logger
	logLevel string
	client   *mongo.Client
	*MongoOpt
}

func NewMongoDB(id, prefix string) *mongoDB {
	return &mongoDB{
		MongoOpt: new(MongoOpt),
		id:       id,
		prefix:   prefix,
	}
}

func (mdb *mongoDB) ID() string {
	return mdb.id
}

func (mdb *mongoDB) InitFlags() {
	prefix := mdb.prefix

	if prefix != "" {
		prefix += "-"
	}

	flag.StringVar(
		&mdb.uri,
		fmt.Sprintf("%smongo-uri", prefix),
		"",
		"MongoDB connection URI",
	)

	flag.StringVar(
		&mdb.database,
		fmt.Sprintf("%smongo-database", prefix),
		"",
		"MongoDB database name",
	)

	flag.Uint64Var(
		&mdb.maxPoolSize,
		fmt.Sprintf("%smongo-max-pool-size", prefix),
		100,
		"Maximum number of connections in the connection pool - Default 100",
	)

	flag.Uint64Var(
		&mdb.minPoolSize,
		fmt.Sprintf("%smongo-min-pool-size", prefix),
		10,
		"Minimum number of connections in the connection pool - Default 10",
	)

	flag.IntVar(
		&mdb.maxConnIdleTime,
		fmt.Sprintf("%smongo-max-conn-idle-time", prefix),
		300,
		"Maximum amount of time a connection can remain idle in seconds - Default 300",
	)

	flag.IntVar(
		&mdb.connectTimeout,
		fmt.Sprintf("%smongo-connect-timeout", prefix),
		10,
		"Connection timeout in seconds - Default 10",
	)

	flag.IntVar(
		&mdb.serverSelectionTimeout,
		fmt.Sprintf("%smongo-server-selection-timeout", prefix),
		30,
		"Server selection timeout in seconds - Default 30",
	)
}

func (mdb *mongoDB) Activate(serviceCtx sctx.ServiceContext) error {
	mdb.logger = serviceCtx.Logger(mdb.id)
	mdb.logLevel = serviceCtx.LogLevel()

	mdb.logger.Infof(
		"mongodb initialized (uri=%s, database=%s, log_level=%s)",
		mdb.uri,
		mdb.database,
		mdb.logLevel,
	)

	if mdb.database == "" {
		return errors.New("mongodb database name is required")
	}

	client, err := mdb.connect()
	if err != nil {
		return err
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return err
	}

	mdb.client = client
	mdb.logger.Info("connected to mongodb successfully")

	return nil
}

func (mdb *mongoDB) Stop() error {
	if mdb.client == nil {
		return nil
	}

	mdb.logger.Info("closing mongodb connection...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return mdb.client.Disconnect(ctx)
}

func (mdb *mongoDB) GetClient() *mongo.Client {
	return mdb.client
}

func (mdb *mongoDB) GetDatabase() *mongo.Database {
	return mdb.client.Database(mdb.database)
}

func (mdb *mongoDB) GetCollection(collectionName string) *mongo.Collection {
	return mdb.GetDatabase().Collection(collectionName)
}

func (mdb *mongoDB) connect() (*mongo.Client, error) {
	clientOptions := options.Client().
		ApplyURI(mdb.uri).
		SetMaxPoolSize(mdb.maxPoolSize).
		SetMinPoolSize(mdb.minPoolSize).
		SetMaxConnIdleTime(time.Duration(mdb.maxConnIdleTime) * time.Second).
		SetConnectTimeout(time.Duration(mdb.connectTimeout) * time.Second).
		SetServerSelectionTimeout(time.Duration(mdb.serverSelectionTimeout) * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}
