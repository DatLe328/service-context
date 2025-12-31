package mongoc

import (
	"context"
	"os"
	"testing"

	sctx "github.com/DatLe328/service-context"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var testServiceCtx sctx.ServiceContext

func TestMain(m *testing.M) {
	testServiceCtx = sctx.NewServiceContext(
		sctx.WithName("test"),
		sctx.WithComponent(NewMongoDB("mongodb", "")),
	)

	if err := testServiceCtx.Load(); err != nil {
		panic(err)
	}

	defer testServiceCtx.Stop()

	code := m.Run()
	os.Exit(code)
}

func TestMongoDB_Connect_Success(t *testing.T) {
	mongoComp := testServiceCtx.MustGet("mongodb").(*mongoDB)

	err := mongoComp.client.Ping(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMongoDB_InsertFind_Success(t *testing.T) {
	mongoComp := testServiceCtx.MustGet("mongodb").(*mongoDB)
	col := mongoComp.GetCollection("users")

	ctx := context.Background()

	_, err := col.InsertOne(ctx, bson.M{"name": "dat"})
	if err != nil {
		t.Fatal(err)
	}

	err = col.FindOne(ctx, bson.M{"name": "dat"}).Err()
	if err != nil {
		t.Fatal(err)
	}
}
