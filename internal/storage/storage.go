package storage

import "gorm.io/gorm"

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

// создание записи
func (s *Storage) Create(track *LikedTrack) error {
	return s.db.Create(track).Error
}

// удаление по айди записи
func (s *Storage) Delete(id int64) error {
	return s.db.Delete(&LikedTrack{}, id).Error
}

// получить все айди в бд
func (s *Storage) GetAllID() ([]int64, error) {
	var ids []int64
	err := s.db.Model(&LikedTrack{}).Pluck("id", &ids).Error
	return ids, err
}

// получить трек по айди
func (s *Storage) GetByID(id int64) (*LikedTrack, error) {
	var track LikedTrack
	err := s.db.First(&track, id).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}
