package activity

import "go.mongodb.org/mongo-driver/mongo"

func getActivityCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("activity")
}
