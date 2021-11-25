package entities

type Battle struct {
    Base
    Finished bool
    BattleParticipants []BattleParticipant
}

type BattleActionType string
const (
    BattleActionTypeWeapon BattleActionType = "weapon"
    BattleActionTypeSpell = "spell"
    BattleActionTypeItem = "item"
    BattleActionTypeEscape = "escape"
)

type BattleAction struct {
    BattleID string
    CauserID string
    TargetID string
    ActionType BattleActionType
    Value int64
}

type BattleParticipantType string
const (
    BattleParticipantTypeCreature BattleParticipantType = "creature"
    BattleParticipantTypePlayer = "player"
)

type BattleParticipant struct {
    Base
    BattleID string
    ParticipantType BattleParticipantType
    ParticipantID string
}




