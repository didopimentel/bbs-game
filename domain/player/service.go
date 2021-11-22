package player

import (
    "bbs-game/domain/entities"
    "gorm.io/gorm"
)

type Service struct {
    db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{
        db: db,
    }
}

func (s *Service) Create(player entities.Player) (entities.Player, error) {
    player.HP = 100
    player.TotalHP = 100
    player.Experience = 0
    player.Level = 1
    tx := s.db.Create(&player)

    if tx.Error != nil {
        return entities.Player{}, tx.Error
    }

    return player, nil
}

func (s *Service) Get(id string) (entities.Player, error) {
    player := entities.Player{}
    tx := s.db.First(&player, id)

    if tx.Error != nil {
        return entities.Player{}, tx.Error
    }

    return player, nil
}

func (s *Service) FindActivityForPlayer(playerID string) (entities.Player, error) {
    player := entities.Player{}
    tx := s.db.First(&player, playerID)

    if tx.Error != nil {
        return entities.Player{}, tx.Error
    }



    return player, nil
}