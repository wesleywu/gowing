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
	col := client.Collection("video_collection")
	//filter := bson.D{}
	filter := bson.M{
		"contentType": bson.M{
			"$in": bson.A{int32(1), int32(2)},
		},
	}
	filter1 := bson.D{
		{"contentType", bson.M{
			"$in": bson.A{int32(1), int32(2)},
		}},
	}
	count, err := col.CountDocuments(ctx, filter1)
	assert.Equal(t, count, int64(2))
	cur, err := col.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	cur.All(ctx, nil)
	lifespan.OnShutdown(ctx)
}
