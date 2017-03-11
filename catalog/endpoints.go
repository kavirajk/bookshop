package catalog

import (
	"net/http"

	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints combine all the catalog service endpoints under single type.
type Endpoints struct {
	SearchEndpoint endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
}

// MakeEndpoints returns Endpoints type which is the combination of
// all the catalog service endpoints.
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		SearchEndpoint: MakeSearchEndpoint(s),
		GetEndpoint:    MakeGetEndpoint(s),
	}
}

func MakeSearchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(searchRequest)
		books, e := s.Search(ctx, req.Q)
		if e != nil {
			return searchResponse{Books: make([]Book, 0), Error: e}, nil
		}
		return searchResponse{Books: books, Status: http.StatusOK}, nil
	}
}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		book, e := s.Get(ctx, req.ID)
		if e != nil {
			return getResponse{Book: nil, Error: e}, nil
		}
		return getResponse{Book: &book}, nil
	}
}

type searchRequest struct {
	Q string `json:"q"`
}

type searchResponse struct {
	Status int    `json:"-"`
	Books  []Book `json:"books,omitempty"`
	Error  error  `json:"error,omitempty"`
}

func (r searchResponse) status() int {
	return r.Status
}

func (r searchResponse) error() error {
	return r.Error
}

type getRequest struct {
	ID string `json:"id"`
}

type getResponse struct {
	Status int   `json:"-"`
	Book   *Book `json:"book,omitempty"`
	Error  error `json:"error,omitempty"`
}

func (l getResponse) status() int {
	return l.Status
}

func (r getResponse) error() error {
	return r.Error
}
