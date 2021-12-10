package repositories

import (
	"github.com/go-pg/pg/v10"
	"market/models"
)

// Products - products repository interface
type Products interface {
	GetProducts() ([]*models.Product, error)
}

type products struct {
	db *pg.DB
}

// NewProducts - new products repository
func NewProducts(db *pg.DB) Products {
	return &products{
		db: db,
	}
}

// GetProducts - get all products from database
func (p *products) GetProducts() ([]*models.Product, error) {
	var products []*models.Product
	err := p.db.Model(&products).Select()
	if err != nil {
		return nil, err
	}

	return products, nil
}
