package storage

import "gorm.io/gorm"

type storage struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *storage {
	return &storage{DB: db}
}

func (s *storage) BeginTransaction() *gorm.DB {
	return s.DB.Begin()
}
