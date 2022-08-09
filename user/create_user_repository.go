package user

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(db *mongo.Database) func(context.Context, User) error {
	return func(ctx context.Context, user User) error {
		collection := getUserCollection(db)
		_, err := collection.InsertOne(ctx, user)
		return err
	}
}
