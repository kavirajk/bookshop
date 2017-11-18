package order

// Repo abstracts all the persistant storage operations of Order Service
type Repo interface {
	Create(order *Order) error
	Save(order *Order) error
	GetByID(ID string) (Order, error)
	ListByUser(userID string) ([]Order, error)
	Drop() error
}
