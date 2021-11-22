package item

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

func (s *Service) Create(item entities.Item) (entities.Item, error) {
    tx := s.db.Create(&item)

    if tx.Error != nil {
        return entities.Item{}, tx.Error
    }

    return item, nil
}

func (s *Service) Get(id string) (entities.Item, error) {
    item := entities.Item{}
    tx := s.db.First(&item, id)

    if tx.Error != nil {
        return entities.Item{}, tx.Error
    }

    return item, nil
}