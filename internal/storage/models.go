package storage

import "time"

// структура одного лайкнутого трека
// для sqlite структура та же
type LikedTrack struct {
	ID       int64 `gorm:"primaryKey"`
	Title    string
	Artist   string
	FilePath string // полный путь файла на диске
	LikedAt  time.Time
}

type UpdateLog struct {
	ID         int64 `gorm:"primaryKey"`
	CheckedAt  time.Time
	Tool       string
	OldVersion string
	NewVersion string
	Updated    bool
}
