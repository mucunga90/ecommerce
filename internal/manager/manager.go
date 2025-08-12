package manager

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/mucunga90/ecommerce/internal"
	"github.com/mucunga90/ecommerce/internal/api"
)

type manager struct {
	storage    storage
	publisher  publisher
	adminEmail string
}

func New(adminEmail string, storage storage, publisher publisher) *manager {
	return &manager{adminEmail: adminEmail, storage: storage, publisher: publisher}
}

func (m *manager) ProductUpload(items []api.ProductPayload) error {
	for _, item := range items {
		categoryID, err := m.ensureCategory(item.Categories)
		if err != nil {
			return err
		}

		product := internal.Product{
			Name:       item.Name,
			CategoryID: categoryID,
			Price:      item.Price,
		}

		if err := m.storage.CreateProduct(&product); err != nil {
			return err
		}
	}

	return nil
}

func (h *manager) ensureCategory(categories []string) (uuid.UUID, error) {
	var parentID uuid.UUID

	for _, name := range categories {
		storedCat, err := h.storage.GetCategoryByName(name, &parentID)
		if err != nil {
			return uuid.Nil, err
		}

		if storedCat != nil {
			parentID = storedCat.ID
			continue
		}

		cat := &internal.Category{
			Name:     name,
			ParentID: &parentID,
		}

		if err := h.storage.CreateCategory(cat); err != nil {
			return uuid.Nil, err
		}
		parentID = cat.ID
	}

	return parentID, nil
}

func (m *manager) ProductAveragePrice(categoryName string) (float64, error) {
	storedCat, err := m.storage.GetCategory(categoryName)
	if err != nil {
		return 0, fmt.Errorf("failed to get category: %w", err)
	}
	avgPrice, err := m.storage.GetAveragePriceForCategory(storedCat.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to get average price for category %s: %w", categoryName, err)
	}
	return avgPrice, nil
}

func (m *manager) CreateOrder(o *internal.Order) error {
	var total float64
	for i := range o.Items {
		total += o.Items[i].UnitPrice * float64(o.Items[i].Quantity)
	}
	o.Total = total

	ctx := context.Background()
	tx := m.storage.BeginTransaction()

	// Save order + items
	if err := tx.Create(&o).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Publish event
	event := internal.OrderCreatedEvent{
		OrderID:      o.ID,
		CustomerName: o.Customer.Name,
		AdminEmail:   m.adminEmail,
		Items:        o.Items,
		Total:        o.Total,
	}

	if err := m.publisher.Publish(ctx, "order.created", event); err != nil {
		tx.Rollback() // rollback DB if publish fails
		return fmt.Errorf("failed to publish event: %w", err)
	}

	// Commit transaction only if publish succeeded
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return m.storage.CreateOrder(o)
}

type storage interface {
	BeginTransaction() *gorm.DB

	CreateCategory(c *internal.Category) error

	CreateProduct(p *internal.Product) error
	GetCategory(name string) (*internal.Category, error)
	GetCategoryByName(name string, parentID *uuid.UUID) (*internal.Category, error)
	GetAveragePriceForCategory(categoryID uuid.UUID) (float64, error)
	CreateOrder(o *internal.Order) error
}

type publisher interface {
	Publish(ctx context.Context, topic string, message interface{}) error
}
