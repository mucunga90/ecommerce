package storage

import (
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/mucunga90/ecommerce/internal"
)

func (s *storage) CreateCategory(c *internal.Category) error {
	if err := s.DB.Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) GetCategory(name string) (*internal.Category, error) {
	var c internal.Category
	if err := s.DB.Where("name = ?", name).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (s *storage) GetCategoryByName(name string, parentID *uuid.UUID) (*internal.Category, error) {
	var c internal.Category
	if err := s.DB.Where("name = ? AND parent_id IS NOT DISTINCT FROM ?", name, parentID).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (s *storage) GetCategoryTree(categoryID uuid.UUID) ([]internal.Category, error) {
	var result []internal.Category

	query := `
        WITH RECURSIVE category_tree AS (
            SELECT id, name, parent_id
            FROM categories
            WHERE id = ? -- starting category id

            UNION ALL

            SELECT c.id, c.name, c.parent_id
            FROM categories c
            INNER JOIN category_tree ct ON c.parent_id = ct.id
        )
        SELECT * FROM category_tree;
    `

	if err := s.DB.Raw(query, categoryID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
