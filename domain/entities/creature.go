package entities

type Creature struct {
    Base
    Name string `json:"name"`
    Level int64 `json:"level"`
    Damage string `json:"damage"`
    Experience int64 `json:"experience"`
    HP int64 `json:"hp"`
}
