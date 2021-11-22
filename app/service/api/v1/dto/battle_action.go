package dto

// TODO: remove causer id to take from token AND amount to take from a calculated number on backend
type CreateActionInput struct {
    CauserID string `json:"causer_id"`
    TargetID string `json:"target_id"`
    ActionType BattleActionType `json:"action_type"`
}

type BattleActionType string
const (
    BattleActionTypeWeapon BattleActionType = "weapon"
    BattleActionTypeSpell = "spell"
    BattleActionTypeItem = "item"
    BattleActionTypeEscape = "escape"
)

type GenerateNextRound struct {
    BattleID string `json:"battle_id"`
}
