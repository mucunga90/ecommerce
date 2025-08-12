package internal

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:text;not null"`
	Email     string    `gorm:"type:text;uniqueIndex;not null"`
	Phone     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"default:now()"`
}
