package storage

import (
	"github.com/mucunga90/ecommerce/internal"
)

func (s *storage) CreateOrder(o *internal.Order) error {
	if err := s.DB.Create(o).Error; err != nil {
		return err
	}
	return nil

}
