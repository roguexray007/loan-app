package base

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	AttributeID = "id"

	// Random integer ceil value
	maxRandomIntCeil int64 = 1111111111111
)

type Model struct {
	ID        int64 `json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type IModel interface {
	TableName() string
	EntityName() string
	GetID() int64
}

// --------------- Getters -------------------- //

func (m *Model) GetID() int64 {
	return m.ID
}

func (m *Model) GetCreatedAt() int64 {
	return m.CreatedAt
}

func (m *Model) GetUpdatedAt() int64 {
	return m.UpdatedAt
}

// --------------- End Getters ------------------ //

// BeforeCreate sets new id for the model.
func (m *Model) BeforeCreate(scope *gorm.Scope) {
	idField, _ := scope.FieldByName("ID")

	if idField.IsBlank {
		v, _ := New()
		_ = idField.Set(v)
	}
}

func New() (int64, error) {
	nanotime := time.Now().UnixNano()

	random := rand.Int63n(maxRandomIntCeil)

	id := nanotime + random

	return id, nil
}

type IEntity interface {
	GetID() int64
}
