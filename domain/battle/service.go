package battle

import (
    "bbs-game/domain/entities"
    "errors"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

var ErrBattleAlreadyFinished = errors.New("battle has already finished")
var ErrPlayerNotInBattle = errors.New("player is not in this battle")

type Service struct {
    db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{
        db: db,
    }
}

var dices = map[int][]int{
    6: {1, 2, 3, 4, 5, 6},
    12: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
}


type CreateInput struct {
    BattleParticipants []entities.BattleParticipant `json:"battle_participants"`
}
type CreateOutput struct {
    Battle entities.Battle
}
func (s *Service) Create(input CreateInput) (CreateOutput, error) {
    participants := make([]entities.BattleParticipant, 0)
    var creatureID string
    for _, bp := range input.BattleParticipants {
        if bp.ParticipantType == entities.BattleParticipantTypeCreature {
            creatureID = bp.ParticipantID
        } else {
            // only append player participants now because creatures will be appended later
            participants = append(participants, entities.BattleParticipant{
                ParticipantType: bp.ParticipantType,
                ParticipantID:   bp.ParticipantID,
            })
        }
    }

    battle := entities.Battle{}
    err := s.db.Transaction(func(tx *gorm.DB) error {
        var creature entities.Creature
        tx = s.db.Table("creatures").First(&creature, "id = ?", creatureID)
        if tx.Error != nil {
            return tx.Error
        }

        creature.ID = ""
        tx = s.db.Table("battle_creatures").Create(&creature)
        if tx.Error != nil {
            return tx.Error
        }

        // append creature participant
        participants = append(participants, entities.BattleParticipant{
            ParticipantType: entities.BattleParticipantTypeCreature,
            ParticipantID:   creature.ID,
        })

        battle.BattleParticipants = participants
        tx = s.db.Table("battles").Create(&battle)
        if tx.Error != nil {
            return tx.Error
        }

        return nil
    })

    if err != nil {
        return CreateOutput{}, err
    }

    return CreateOutput{
        Battle: battle,
    }, nil
}

type CreateActionInput struct {
    BattleID string
    CauserID string
    TargetID string
    ActionType entities.BattleActionType
}

//TODO: delete
func (s *Service) CreateAction(input CreateActionInput) (entities.BattleAction, error) {
    var battle entities.Battle
    tx := s.db.First(&battle, input.BattleID)
    if tx.Error != nil {
        return entities.BattleAction{}, tx.Error
    }

    causerParticipant := entities.BattleParticipant{}
    targetParticipant := entities.BattleParticipant{}
    for _, p := range battle.BattleParticipants {
        if p.ID == input.CauserID {
            causerParticipant = p
        }
        if p.ID == input.TargetID {
            targetParticipant = p
        }
    }

    var damage int64
    if causerParticipant.ParticipantType == entities.BattleParticipantTypePlayer {
        var player entities.Player
        tx = s.db.Table("players").First(&player, causerParticipant.ParticipantID)
        if tx.Error != nil {
            return entities.BattleAction{}, tx.Error
        }

        totalDamage, err := getDamage(player.Damage)
        if err != nil {
            return entities.BattleAction{}, err
        }
        damage = int64(totalDamage)

        var creature entities.Creature
        tx = s.db.Table("creatures").First(&creature, targetParticipant.ParticipantID)
        if tx.Error != nil {
            return entities.BattleAction{}, tx.Error
        }

        creature.HP = creature.HP - damage
        tx = s.db.Save(&creature)
        if tx.Error != nil {
            return entities.BattleAction{}, tx.Error
        }
    }
    if causerParticipant.ParticipantType == entities.BattleParticipantTypeCreature {
        var creature entities.Creature
        tx = s.db.Table("creatures").First(&creature, causerParticipant.ParticipantID)
        if tx.Error != nil {
            return entities.BattleAction{}, tx.Error
        }

        totalDamage, err := getDamage(creature.Damage)
        if err != nil {
            return entities.BattleAction{}, err
        }

        damage = int64(totalDamage)

        var player entities.Player
        tx = s.db.Table("players").First(&player, targetParticipant.ParticipantID)
        if tx.Error != nil {
            return entities.BattleAction{}, tx.Error
        }

        player.HP = player.HP - damage
        tx = s.db.Save(&player)
        if tx.Error != nil {
            return entities.BattleAction{}, tx.Error
        }
    }

    battleAction := entities.BattleAction{
        BattleID:   input.BattleID,
        CauserID:   input.CauserID,
        TargetID:   input.TargetID,
        ActionType: input.ActionType,
        Value:      damage,
    }

    err := s.db.Transaction(func(tx *gorm.DB) error {
        tx = tx.Table("battle_actions").Create(&battleAction)
        if tx.Error != nil {
            return tx.Error
        }
        return nil
    })
    if err != nil {
        return entities.BattleAction{}, s.db.Error
    }

    return battleAction, nil
}

type GenerateNextRoundInput struct {
    BattleID string
    PlayerID string
}
type GenerateNextRoundOutput struct {
    BattleActions []entities.BattleAction
    PlayerDied bool
    CreatureDied bool
    ExperienceGained int64
    GainedLevel bool
}
func (s *Service) GenerateNextRound(input GenerateNextRoundInput) (GenerateNextRoundOutput, error) {
    actions := make([]entities.BattleAction, 0)
    playerDied, creatureDied, gainedLevel := false, false, false
    var gainedExperience int64
    var battle entities.Battle
    tx := s.db.Preload(clause.Associations).Table("battles").First(&battle, "id = ?", input.BattleID)
    if tx.Error != nil {
        return GenerateNextRoundOutput{}, tx.Error
    }

    if battle.Finished {
        return GenerateNextRoundOutput{}, ErrBattleAlreadyFinished
    }

    var player entities.Player
    var creature entities.Creature
    var playerParticipant, creatureParticipant entities.BattleParticipant

    for _, bp := range battle.BattleParticipants {
        switch bp.ParticipantType {
        case entities.BattleParticipantTypePlayer:
            tx = s.db.Table("players").First(&player, "id = ?", bp.ParticipantID)
            if tx.Error != nil {
                return GenerateNextRoundOutput{}, tx.Error
            }

            if player.ID != input.PlayerID {
                return GenerateNextRoundOutput{}, ErrPlayerNotInBattle
            }
            playerParticipant = bp
        case entities.BattleParticipantTypeCreature:
            tx = s.db.Table("battle_creatures").First(&creature, "id = ?", bp.ParticipantID)
            if tx.Error != nil {
                return GenerateNextRoundOutput{}, tx.Error
            }
            creatureParticipant = bp
        }
    }


    playerDamage, err := getDamage(player.Damage)
    if err != nil {
        return GenerateNextRoundOutput{}, err
    }
    creatureDamage, err := getDamage(creature.Damage)
    if err != nil {
        return GenerateNextRoundOutput{}, err
    }

    s.db.Transaction(func(tx *gorm.DB) error {
        creature.HP = creature.HP - int64(playerDamage)
        log1 := entities.BattleAction{
            BattleID:   battle.ID,
            CauserID:   playerParticipant.ID,
            TargetID:   creatureParticipant.ID,
            ActionType: entities.BattleActionTypeWeapon,
            Value:      int64(playerDamage),
        }
        player.HP = player.HP - int64(creatureDamage)
        log2 := entities.BattleAction{
            BattleID:   battle.ID,
            CauserID:   creatureParticipant.ID,
            TargetID:   playerParticipant.ID,
            ActionType: entities.BattleActionTypeWeapon,
            Value:      int64(creatureDamage),
        }

        actions = append(actions, log1, log2)

        if player.HP <= 0 {
            playerDied = true
            return nil
        }

        if creature.HP <= 0 {
            gainedExperience = creature.Experience
            creatureDied = true

            // Restore HP
            player.HP = player.TotalHP
            player.Experience = player.Experience + gainedExperience

            nextLevelExperience := entities.LevelExperience[player.Level + 1]

            if player.Experience > nextLevelExperience {
                player.Level++
                gainedLevel = true
            }

            tx = s.db.Table("battle_creatures").Delete(&creature)
            if tx.Error != nil {
                return tx.Error
            }
        }

        tx = s.db.Table("players").Save(&player)
        if tx.Error != nil {
            return tx.Error
        }

        tx = s.db.Table("battle_creatures").Save(&creature)
        if tx.Error != nil {
            return tx.Error
        }

        tx = s.db.Table("battle_actions").Create(&actions)
        if tx.Error != nil {
            return tx.Error
        }

        return nil
    })

    return GenerateNextRoundOutput{
        BattleActions: actions,
        PlayerDied:    playerDied,
        CreatureDied:  creatureDied,
        ExperienceGained: gainedExperience,
        GainedLevel: gainedLevel,
    }, nil
}

func getDamage(d string) (int, error) {
    totalDamage := 0

    diceCount, diceSize, additionalDamage, err := getDiceDamage(d)
    if err != nil {
        return 0, err
    }

    for i := diceCount; i > 0; i-- {
        rand.Seed(time.Now().UnixNano())
        totalDamage = totalDamage + dices[diceSize][rand.Intn(len(dices[diceSize]))]
    }
    return totalDamage + additionalDamage, nil
}

func getDiceDamage(d string) (int, int, int, error) {
    splitted1 := strings.Split(d, "d")
    diceCount, err := strconv.Atoi(splitted1[0])
    if err != nil {
        return 0, 0, 0, err
    }
    splitted2 := strings.Split(splitted1[1], "+")
    diceSize, err := strconv.Atoi(splitted2[0])
    if err != nil {
        return 0, 0, 0, err
    }
    if len(splitted2) == 1 {
        return diceCount, diceSize, 0, nil
    }

    additionalDamage, err := strconv.Atoi(splitted2[1])
    if err != nil {
        return 0, 0, 0, err
    }

    return diceCount, diceSize, additionalDamage, nil
}

func (s *Service) Get(id string) (entities.Battle, error) {
    battle := entities.Battle{}
    tx := s.db.First(&battle, "id = ?", id)

    if tx.Error != nil {
        return entities.Battle{}, tx.Error
    }

    return battle, nil
}

func (s *Service) ListActionsByBattle(id string) ([]entities.BattleAction, error) {
    battleActions := make([]entities.BattleAction, 0)

    tx := s.db.Find(&battleActions)
    if tx.Error != nil {
        return battleActions, tx.Error
    }

    return battleActions, nil
}