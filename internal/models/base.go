package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Base struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	SoftDeletedAt *time.Time         `bson:"soft_deleted_at,omitempty" json:"soft_deleted_at,omitempty"`
}

// before create hook to set a uuid
func (b *Base) BeforeCreate() {
	now := time.Now()

	if b.ID.IsZero() {
		b.ID = primitive.NewObjectID()
	}

	b.CreatedAt = now
	b.UpdatedAt = now
}

// before update will update the UpdatedAt timestamo

func (b *Base) BeforeUpdate() {
	b.UpdatedAt = time.Now()
}

// after delete will delete deleted at timestamp

func (b *Base) BeforeSoftDelete() {
	now := time.Now()
	// soft deleted at can point to a nil pointer. hence the weird pointer assignment
	b.SoftDeletedAt = &now
	b.UpdatedAt = time.Now()
}
