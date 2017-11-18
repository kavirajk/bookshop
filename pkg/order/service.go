package order

import (
	"context"
	"errors"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Service interface {
	// PlaceOrder creates an order for particular book.
	PlaceOrder(ctx context.Context, bookID string) (Order, error)

	// GetUserOrders returns list of orders placed by an user.
	GetUserOrders(ctx context.Context, userID string) ([]Order, error)

	// CancelOrder cancels the particular order of an user.
	CancelOrder(ctx context.Context, userID string, orderID string) error
}

type basicService struct {
	r Repo
}

// NewOrderService return basic Service implementation.
func NewService(r Repo) Service {
	return basicService{r}
}

// PlaceOrder creates an order for particular book.
func (s basicService) PlaceOrder(ctx context.Context, bookID string) (Order, error) {
	return Order{}, nil
}

// GetUserOrders return all the orders placed by particular user.
func (s basicService) GetUserOrders(ctx context.Context, userID string) ([]Order, error) {
	return nil, nil
}

// CancelOrder cancels a particular order placed by particular user.
// If particular is not placed by the user, it returns ErrOrderNotFound.
func (s basicService) CancelOrder(ctx context.Context, userID string, orderID string) error {
	return nil
}

type Middleware func(Service) Service
