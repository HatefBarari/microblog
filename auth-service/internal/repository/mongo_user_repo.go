package repository

import (
	"context"
	"errors"

	"github.com/HatefBarari/microblog-auth/internal/domain"
	"github.com/HatefBarari/microblog-shared/pkg/mongo" // ← پکیج داخلی ما
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	driver "go.mongodb.org/mongo-driver/mongo" // ← alias برای ErrNoDocuments
)

type mongoUserRepo struct{}

func NewMongoUserRepo() domain.UserRepository {
	return &mongoUserRepo{}
}

func (r *mongoUserRepo) Create(ctx context.Context, u *domain.User) error {
	res, err := mongo.UsersColl().InsertOne(ctx, u)
	if err != nil {
		return err
	}
	u.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *mongoUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	res := mongo.UsersColl().FindOne(ctx, bson.M{"email": email})
	if err := res.Err(); err != nil {
		if errors.Is(err, driver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	var u domain.User
	if err := res.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *mongoUserRepo) UpdateVerified(ctx context.Context, userID string, verified bool) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = mongo.UsersColl().UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"verified": verified}})
	return err
}