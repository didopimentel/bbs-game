package persistence

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewDB(addr string) (*gorm.DB, error) {
    return gorm.Open(postgres.Open(addr), &gorm.Config{})
}
