package order

import "context"

type Service interface {
	PlaceOrder(ctx context.Context, bookID string) (Order, error)
	CancelOrder(ctx context.Context, orderID string) error
}
