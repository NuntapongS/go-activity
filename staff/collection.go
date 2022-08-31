package staff

import "go.mongodb.org/mongo-driver/mongo"

func getStaffCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection("staffs")
}
