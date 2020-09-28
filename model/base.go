package model

import "time"

type Collection string

type Base struct {
	ID        string    `firestore:"id"`
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

func (b *Base) GetID() string {
	return b.ID
}

func (b *Base) SetID(id string) {
	b.ID = id
}

func (b *Base) GetCreatedAt() time.Time {
	return b.CreatedAt
}

func (b *Base) SetCreatedAt(createdAt time.Time) {
	b.CreatedAt = createdAt
}

func (b *Base) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}

func (b *Base) SetUpdatedAt(updatedAt time.Time) {
	b.UpdatedAt = updatedAt
}
