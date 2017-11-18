package order

import (
	"net/http"

	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints combine all the order service endpoints under single type.
type Endpoints struct {
	PlaceOrderEndpoint    endpoint.Endpoint
	GetUserOrdersEndpoint endpoint.Endpoint
	CancelOrderEndpoint   endpoint.Endpoint
}

// MakeEndpoints returns Endpoints type which is the combination of
// all the order service endpoints.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		PlaceOrderEndpoint:    MakePlaceOrderEndpoint(s),
		GetUserOrdersEndpoint: MakeGetUserOdersEndpoint(s),
		CancelOrderEndpoint:   MakeCancelOrderEndpoint(s),
	}
}

func MakePlaceOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(placeOrderRequest)
		order, e := s.PlaceOrder(ctx, req.BookID)
		if e != nil {
			return placeOrderResponse{Order: nil, Error: e}, nil
		}
		return placeOrderResponse{Order: &order, Status: http.StatusCreated}, nil
	}
}

func MakeGetUserOdersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserOrdersRequest)
		orders, e := s.GetUserOrders(ctx, req.UserID)
		if e != nil {
			return getUserOrdersResponse{Orders: make([]Order, 0), Error: e}, nil
		}
		return getUserOrdersResponse{Orders: orders}, nil
	}
}

func MakeCancelOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(cancelOrderRequest)
		e := s.CancelOrder(ctx, req.UserID, req.OrderID)
		if e != nil {
			return cancelOrderResponse{Error: e}, nil
		}
		return cancelOrderResponse{nil}, nil
	}
}

type placeOrderRequest struct {
	BookID string `json:"book_id"`
}

type placeOrderResponse struct {
	Status int    `json:"-"`
	Order  *Order `json:"order,omitempty"`
	Error  error  `json:"error,omitempty"`
}

func (r placeOrderResponse) status() int {
	return r.Status
}

func (r placeOrderResponse) error() error {
	return r.Error
}

type getUserOrdersRequest struct {
	UserID string `json:"user_id"`
}

type getUserOrdersResponse struct {
	Orders []Order `json:"orders,omitempty"`
	Error  error   `json:"error,omitempty"`
}

func (r getUserOrdersResponse) error() error {
	return r.Error
}

type cancelOrderRequest struct {
	UserID  string `json:"user_id"`
	OrderID string `json:"order_id"`
}

type cancelOrderResponse struct {
	Error error `json:"error,omitempty"`
}
