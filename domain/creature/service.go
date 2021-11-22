package creature

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

func (s *Service) Create(creature entities.Creature) (entities.Creature, error) {
    tx := s.db.Create(&creature)

    if tx.Error != nil {
        return entities.Creature{}, tx.Error
    }

    return creature, nil
}

func (s *Service) Get(id string) (entities.Creature, error) {
    creature := entities.Creature{}
    tx := s.db.First(&creature, id)

    if tx.Error != nil {
        return entities.Creature{}, tx.Error
    }

    return creature, nil
}