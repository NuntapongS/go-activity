package staff

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role string

const (
	ADMIN = "ADMIN"
	STAFF = "STAFF"
)

type Staff struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Username       string             `bson:"username" json:"username"`
	HashedPassword string             `bson:"hashed_password" json:"-"`
	Role           Role               `bson:"role" json:"role"`
}
