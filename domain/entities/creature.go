package entities

type Creature struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Level int64 `json:"level"`
    Damage string `json:"damage"`
    Experience int64 `json:"experience"`
    HP int64 `json:"hp"`
}
