package order

type Service interface {
	PlaceOrder(AlbumID string) (Order, error)
	CancelOrder(orderID string) error
}
