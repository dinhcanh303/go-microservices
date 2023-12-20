package domain

import (
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Key          string                 `bson:"key,omitempty,unique" json:"key"`
	SubjectCount uint64                 `bson:"subject_count" json:"subject_count"`
	Subjects     []sharedkernel.Subject `bson:"subjects" json:"subjects"`
	DiObject     sharedkernel.Subject   `bson:"di_object" json:"di_object"`
	InObject     sharedkernel.Subject   `bson:"in_object" json:"in_object"`
	PrObject     sharedkernel.Subject   `bson:"pr_object" json:"pr_object"`
}
