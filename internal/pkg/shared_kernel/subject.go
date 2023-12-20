package sharedkernel

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subject struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Type  string             `bson:"type" json:"type"`
	Image string             `bson:"image" json:"image"`
}
