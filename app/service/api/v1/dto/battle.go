package dto

import "bbs-game/domain/entities"

type BattleCreateInput struct {
    BattleParticipants []BattleParticipant `json:"battle_participants"`
}
type Battle struct {
    ID string `json:"id"`
    RoundOwner string `json:"round_owner"`
}
type BattleParticipant struct {
    ID string `json:"id"`
    BattleID string `json:"battle_id"`
    ParticipantType BattleParticipantType `json:"participant_type"`
    ParticipantID string `json:"participant_id"`
}
type BattleParticipantType string
const (
    BattleParticipantTypeCreature BattleParticipantType = "creature"
    BattleParticipantTypePlayer = "player"
)

type GenerateNextRoundResponse struct {
    BattleActions []entities.BattleAction `json:"battle_actions"`
    PlayerDied bool `json:"player_died"`
    CreatureDied bool `json:"creature_died"`
    ExperienceGained int64 `json:"experience_gained"`
    GainedLevel bool `json:"gained_level"`
}