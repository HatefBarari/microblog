package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

func Connect(uri, dbName string, log *zap.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	if err := c.Ping(ctx, nil); err != nil {
		return err
	}
	client = c
	db = client.Database(dbName)
	log.Info("mongo connected")
	return nil
}

func Client() *mongo.Client        { return client }
func DB() *mongo.Database          { return db }
func UsersColl() *mongo.Collection { return db.Collection("users") }
func ArticlesColl() *mongo.Collection { return db.Collection("articles") }
func CommentsColl() *mongo.Collection { return db.Collection("comments") }
func CategoriesColl() *mongo.Collection { return db.Collection("categories") }
func MediaColl() *mongo.Collection { return db.Collection("media") }