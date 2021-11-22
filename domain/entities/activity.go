package entities

type Activity struct {
    Type ActivityType
    Id string
}

type ActivityType string

const (
    ActivityTypeBattle string = "activity_type_battle"
    ActivityTypeChest string = "activity_type_chest"
)


