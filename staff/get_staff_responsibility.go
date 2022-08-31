package staff

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStaffByUsername(db *mongo.Database) func(context.Context, string) (*Staff, error) {
	return func(ctx context.Context, username string) (*Staff, error) {
		collection := getStaffCollection(db)

		var result Staff

		filter := bson.M{"username": bson.M{"$eq": username}}

		rs := collection.FindOne(ctx, filter)

		err := rs.Decode(&result)

		return &result, err
	}
}
