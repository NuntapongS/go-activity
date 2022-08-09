package user

import (
	"context"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser(db *mongo.Database) func(context.Context) ([]User, error) {
	return func(ctx context.Context) ([]User, error) {
		collection := getUserCollection(db)
		var users []User
		cur, err := collection.Find(ctx, bson.D{})
		for cur.Next(ctx) {
			var elem User
			err := cur.Decode(&elem)
			if err != nil {
				return nil, err
			}
			users = append(users, elem)
		}

		if cur.Err() != nil {
			return nil, cur.Err()
		}

		cur.Close(ctx)
		return users, err

	}
}

func GetUserByID(db *mongo.Database) func(context.Context, string) (*User, error) {
	return func(ctx context.Context, UserID string) (*User, error) {
		collection := getUserCollection(db)
		var results User

		userID, error := strconv.Atoi(UserID)
		if error != nil {
			return nil, error
		}

		filter := bson.M{"user_id": bson.M{"$eq": userID}}

		rs := collection.FindOne(ctx, filter)
		err := rs.Decode(&results)
		if err != nil {
			return nil, err
		}

		return &results, error
	}
}
