package internal

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID uuid.UUID   `gorm:"type:uuid;not null"`
	Customer   Customer    `gorm:"foreignKey:CustomerID"`
	Total      float64     `gorm:"type:numeric(12,2);not null"`
	Status     string      `gorm:"type:text;default:'pending'"`
	CreatedAt  time.Time   `gorm:"default:now()"`
	Shipping   JSONB       `gorm:"type:jsonb"`
	Notes      string      `gorm:"type:text"`
	Items      []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null"`
	Product    Product   `gorm:"foreignKey:ProductID"`
	Quantity   int       `gorm:"not null"`
	UnitPrice  float64   `gorm:"type:numeric(12,2);not null"`
	TotalPrice float64   `gorm:"type:numeric(12,2);not null"`
}

type JSONB map[string]any

func (j JSONB) GormDataType() string {
	return "jsonb"
}
