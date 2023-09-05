package mongodb

import (
	"context"
	"fmt"
	"sync"

	"github.com/WesleyWu/gowing/util/gworm/mongodb/codecs"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wesleywu/go-lifespan/lifespan"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const HookName = "mongodb"

var (
	mongodbDatabase   *mongo.Database
	mongodbClient     *mongo.Client
	mongodbClientOnce sync.Once
	bootstrapErr      error
)

func init() {
	lifespan.RegisterBootstrapHook(HookName, true, setupConnection)
	lifespan.RegisterShutdownHook(HookName, false, closeConnection)
}

func setupConnection(ctx context.Context) error {
	mongodbClientOnce.Do(func() {
		host := getMongoDBHost(ctx)
		port := getMongoDBPort(ctx)
		username := getMongoDBUsername(ctx)
		password := getMongoDBPassword(ctx)
		databaseName := getMongoDBDatabase(ctx)

		uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", username, password, host, port, databaseName)
		g.Log().Infof(ctx, "connecting to mongodb uri: `%s`", uri)
		reg := codecs.Register(bson.NewRegistryBuilder()).Build()
		clientOpts := &options.ClientOptions{
			Registry: reg,
		}
		if mongodbClient, bootstrapErr = mongo.Connect(ctx, options.Client().ApplyURI(uri), clientOpts); bootstrapErr != nil {
			return
		}
		if bootstrapErr = mongodbClient.Ping(ctx, readpref.Primary()); bootstrapErr != nil {
			return
		}
		mongodbDatabase = mongodbClient.Database(databaseName)
	})
	return bootstrapErr
}

func ensureConnection(ctx context.Context) {
	if mongodbDatabase == nil {
		g.Log().Warningf(ctx, "Mongodb connection was not successfully set up during lifespan bootstrap phase, now setting up...")
		err := setupConnection(ctx)
		if err != nil {
			g.Log().Fatalf(ctx, "Mongodb connection error: %+v", err)
		}
	}
}

func closeConnection(ctx context.Context) error {
	if mongodbClient != nil {
		return mongodbClient.Disconnect(ctx)
	}
	return nil
}

func GetDatabaseInstance(ctx context.Context) *mongo.Database {
	ensureConnection(ctx)
	return mongodbDatabase
}

func GetCollection(ctx context.Context, collectionName string, opts ...*options.CollectionOptions) *mongo.Collection {
	ensureConnection(ctx)
	return mongodbDatabase.Collection(collectionName)
}

func TestConnectivity(ctx context.Context) error {
	ensureConnection(ctx)
	col := mongodbDatabase.Collection("test")
	_, err := col.InsertOne(ctx, bson.D{
		{"foo", "bar"},
	})
	if err != nil {
		return err
	}
	res, err := col.DeleteOne(ctx, bson.D{{"foo", "bar"}})
	if err != nil {
		return err
	}
	if res.DeletedCount != 1 {
		return gerror.New("deleted count should equal to 1")
	}
	return nil
}
