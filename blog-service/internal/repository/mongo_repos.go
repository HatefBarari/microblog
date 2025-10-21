package repository

import (
	"context"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongoDB "github.com/HatefBarari/microblog-shared/pkg/mongo"
)

// ---------- Article ----------
type mongoArticleRepo struct{}

func NewMongoArticleRepo() domain.ArticleRepository { return &mongoArticleRepo{} }

func (r *mongoArticleRepo) Create(ctx context.Context, a *domain.Article) error {
	res, err := mongoDB.DB().Collection("articles").InsertOne(ctx, a)
	if err != nil {
		return err
	}
	a.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *mongoArticleRepo) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var a domain.Article
	if err := mongoDB.DB().Collection("articles").FindOne(ctx, bson.M{"_id": oid}).Decode(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *mongoArticleRepo) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	var a domain.Article
	if err := mongoDB.DB().Collection("articles").FindOne(ctx, bson.M{"slug": slug}).Decode(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *mongoArticleRepo) List(ctx context.Context, filter domain.ListFilter) ([]*domain.Article, int, error) {
	// TODO: build filter, pagination, sort
	cursor, err := mongoDB.DB().Collection("articles").Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}).SetSkip(int64((filter.Page-1)*filter.PageSize)).SetLimit(int64(filter.PageSize)))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var list []*domain.Article
	if err := cursor.All(ctx, &list); err != nil {
		return nil, 0, err
	}
	// TODO: count total
	return list, 100, nil
}

func (r *mongoArticleRepo) Update(ctx context.Context, a *domain.Article) error {
	oid, err := primitive.ObjectIDFromHex(a.ID)
	if err != nil {
		return err
	}
	_, err = mongoDB.DB().Collection("articles").ReplaceOne(ctx, bson.M{"_id": oid, "author_id": a.AuthorID}, a)
	return err
}

func (r *mongoArticleRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mongoDB.DB().Collection("articles").DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *mongoArticleRepo) UpdateStatus(ctx context.Context, id string, status domain.ArticleStatus) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mongoDB.DB().Collection("articles").UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *mongoArticleRepo) UpdateViewCount(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mongoDB.DB().Collection("articles").UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$inc": bson.M{"view_count": 1}})
	return err
}

// ---------- Category ----------
type mongoCategoryRepo struct{}

func NewMongoCategoryRepo() domain.CategoryRepository { return &mongoCategoryRepo{} }

func (r *mongoCategoryRepo) Create(ctx context.Context, c *domain.Category) error {
	res, err := mongoDB.DB().Collection("categories").InsertOne(ctx, c)
	if err != nil {
		return err
	}
	c.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *mongoCategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var c domain.Category
	if err := mongoDB.DB().Collection("categories").FindOne(ctx, bson.M{"_id": oid}).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *mongoCategoryRepo) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var c domain.Category
	if err := mongoDB.DB().Collection("categories").FindOne(ctx, bson.M{"slug": slug}).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *mongoCategoryRepo) ListTree(ctx context.Context) ([]*domain.Category, error) {
	cursor, err := mongoDB.DB().Collection("categories").Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*domain.Category
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	// TODO: build tree
	return list, nil
}

// ---------- Comment ----------
type mongoCommentRepo struct{}

func NewMongoCommentRepo() domain.CommentRepository { return &mongoCommentRepo{} }

func (r *mongoCommentRepo) Create(ctx context.Context, c *domain.Comment) error {
	res, err := mongoDB.DB().Collection("comments").InsertOne(ctx, c)
	if err != nil {
		return err
	}
	c.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *mongoCommentRepo) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var c domain.Comment
	if err := mongoDB.DB().Collection("comments").FindOne(ctx, bson.M{"_id": oid}).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *mongoCommentRepo) ListByArticle(ctx context.Context, articleID string, status domain.CommentStatus) ([]*domain.Comment, error) {
	cursor, err := mongoDB.DB().Collection("comments").Find(ctx, bson.M{"article_id": articleID, "status": status}, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var list []*domain.Comment
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *mongoCommentRepo) UpdateStatus(ctx context.Context, id string, status domain.CommentStatus) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mongoDB.DB().Collection("comments").UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *mongoCommentRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = mongoDB.DB().Collection("comments").DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// ---------- Rating ----------
type mongoRatingRepo struct{}

func NewMongoRatingRepo() domain.RatingRepository { return &mongoRatingRepo{} }

func (r *mongoRatingRepo) Save(ctx context.Context, ra *domain.Rating) error {
	_, err := mongoDB.DB().Collection("ratings").InsertOne(ctx, ra)
	return err
}

func (r *mongoRatingRepo) GetByUserAndTarget(ctx context.Context, userID, targetID, targetType string) (*domain.Rating, error) {
	var ra domain.Rating
	if err := mongoDB.DB().Collection("ratings").FindOne(ctx, bson.M{"user_id": userID, "target_id": targetID, "type": targetType}).Decode(&ra); err != nil {
		return nil, err
	}
	return &ra, nil
}

func (r *mongoRatingRepo) GetAverage(ctx context.Context, targetID, targetType string) (float64, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "target_id", Value: targetID}, {Key: "type", Value: targetType}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "avg", Value: bson.D{{Key: "$avg", Value: "$stars"}}}}}}
	cursor, err := mongoDB.DB().Collection("ratings").Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	var result struct{ Avg float64 }
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Avg, nil
	}
	return 0, nil // no ratings
}

func (r *mongoRatingRepo) Delete(ctx context.Context, userID, targetID, targetType string) error {
	_, err := mongoDB.DB().Collection("ratings").DeleteOne(ctx, bson.M{"user_id": userID, "target_id": targetID, "type": targetType})
	return err
}