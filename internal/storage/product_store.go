package storage

import (
	"github.com/google/uuid"
	"github.com/mucunga90/ecommerce/internal"
)

func (s *storage) CreateProduct(p *internal.Product) error {
	if err := s.DB.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) GetAveragePriceForCategory(categoryID uuid.UUID) (float64, error) {
	var avgPrice float64

	query := `
        WITH RECURSIVE category_tree AS (
            SELECT id, parent_id
            FROM categories
            WHERE id = ?

            UNION ALL

            SELECT c.id, c.parent_id
            FROM categories c
            INNER JOIN category_tree ct ON c.parent_id = ct.id
        )
        SELECT COALESCE(AVG(p.price), 0) 
        FROM products p
        INNER JOIN category_tree ct ON p.category_id = ct.id;
    `

	err := s.DB.Raw(query, categoryID).Scan(&avgPrice).Error
	if err != nil {
		return 0, err
	}

	return avgPrice, nil
}
