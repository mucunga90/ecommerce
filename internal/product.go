package internal

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null"`
	Category    Category  `gorm:"foreignKey:CategoryID"`
	Name        string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text;not null"`
	Price       float64   `gorm:"type:numeric(12,2);not null"`
	CreatedAt   time.Time `gorm:"default:now()"`
}
