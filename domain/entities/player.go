package entities

import (
    "time"
)

type Player struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Damage string `json:"damage"`
    Level int64 `json:"level"`
    Experience int64 `json:"experience"`
    HP int64 `json:"hp"`
    TotalHP int64 `json:"total_hp"`
    AccountID string `json:"account_id"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
