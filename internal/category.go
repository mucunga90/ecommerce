package internal

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ParentID  *uuid.UUID `gorm:"type:uuid"`
	Name      string     `gorm:"type:text;uniqueIndex;not null"`
	CreatedAt time.Time  `gorm:"default:now()"`
}
