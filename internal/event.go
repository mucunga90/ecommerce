package internal

import (
	"github.com/google/uuid"
)

type OrderCreatedEvent struct {
	OrderID       uuid.UUID   `json:"order_id"`
	CustomerPhone string      `json:"customer_phone"`
	CustomerName  string      `json:"customer_name"`
	AdminEmail    string      `json:"admin_email"`
	Items         []OrderItem `json:"items"`
	Total         float64     `json:"total"`
}
