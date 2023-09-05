package mongodb

import (
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/stretchr/testify/assert"
	"github.com/wesleywu/go-lifespan/lifespan"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetDatabaseInstance(t *testing.T) {
	ctx := gctx.New()
	lifespan.OnBootstrap(ctx)

	client := GetDatabaseInstance(ctx)
	assert.NotNil(t, client, "failed to get mongodb client instance")
	c, err := client.Collection("video_collection").Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	c.All(ctx, nil)
	lifespan.OnShutdown(ctx)
}
