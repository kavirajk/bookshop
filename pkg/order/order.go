package order

import (
	"github.com/kavirajk/bookshop/catalog"
	"github.com/kavirajk/bookshop/user"
)

type Order struct {
	ID          string         `json:"id"`
	CreatedBy   *user.User     `json:"created_by"`
	CreatedByID string         `json:"-"`
	Items       []catalog.Book `json:"items" gorm:"many_to_many"`
	TotalPrice  float64        `json:"total_price"`
	Currency    string         `json:"currency"`
}
