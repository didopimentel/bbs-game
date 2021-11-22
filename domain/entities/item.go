package entities

type ItemType string
const (
    ItemTypeSword ItemType = "sword"
    ItemTypeArmor ItemType = "armor"
    ItemTypeScroll ItemType = "scroll"
)

type DamageType string
const (
    DamageTypeNormal ItemType = "normal"
)

type Item struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Type ItemType `json:"type"`
    DamageAmount int64 `json:"damage_amount"`
    HealAmount int64 `json:"heal_amount"`
    DamageType DamageType `json:"damage_type"`
}