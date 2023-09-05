package mongodb

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func getMongoDBUsername(ctx context.Context) string {
	value, err := g.Cfg().GetWithEnv(ctx, "mongodb.username", "gowing")
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return "root"
	} else {
		return value.String()
	}
}

func getMongoDBPassword(ctx context.Context) string {
	value, err := g.Cfg().GetWithEnv(ctx, "mongodb.password", "gowing")
	if err != nil {
		return "example"
	} else {
		return value.String()
	}
}

func getMongoDBHost(ctx context.Context) string {
	value, err := g.Cfg().GetWithEnv(ctx, "mongodb.host", "localhost")
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return "localhost"
	} else {
		return value.String()
	}
}

func getMongoDBPort(ctx context.Context) int {
	value, err := g.Cfg().GetWithEnv(ctx, "mongodb.port", 27017)
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return 27017
	}
	return value.Int()
}

func getMongoDBDatabase(ctx context.Context) string {
	value, err := g.Cfg().GetWithEnv(ctx, "mongodb.database", "gowing")
	if err != nil {
		g.Log().Errorf(ctx, "%+v", err)
		return "test"
	} else {
		return value.String()
	}
}
