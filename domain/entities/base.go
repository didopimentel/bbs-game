package entities

import (
    "github.com/gofrs/uuid"
    "gorm.io/gorm"
    "time"
)

type Base struct {
    ID         string     `sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
    uuid, err := uuid.NewV4()
    if err != nil {
        return err
    }
    base.ID = uuid.String()
    return tx.Error
}
