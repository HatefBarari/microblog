package tests

import (
	"context"
	"testing"

	"github.com/HatefBarari/microblog-shared/pkg/logger"
	"github.com/HatefBarari/microblog-shared/pkg/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoPing(t *testing.T) {
	log, _ := logger.NewFile("debug", "logs/test.log") // ← اسم درست
	assert.NoError(t, mongo.Connect("mongodb://localhost:27017", "testdb", log))
	ctx := context.Background()
	assert.NoError(t, mongo.DB().RunCommand(ctx, bson.M{"ping": 1}).Err())
}