package entities

import (
    "time"
)

type Account struct {
    ID string
    CreatedAt time.Time
    UpdatedAt time.Time
    Email string
    Password string
}
