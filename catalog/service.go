package catalog

import (
	"context"
	"errors"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type Service interface {
	// Search books based on free text
	Search(ctx context.Context, query string) ([]Book, error)

	// List available items based on limit and offset.
	// order takes string in the format "name asc" or "name desc"
	// or in combination of multiple fields like "name asc, isbn desc"
	List(ctx context.Context, order string, limit, offset int) ([]Book, int, error)

	// Get details about single book
	Get(ctx context.Context, id string) (Book, error)
}

type basicService struct {
	r Repo
}

// NewCatalogService return basic Service implementation.
func NewService(r Repo) Service {
	return basicService{r}
}

// Search return books that matches with query.
func (s basicService) Search(ctx context.Context, query string) ([]Book, error) {
	return s.r.Search(query)
}

// Get return a book for the matched ID. Empty book incase of non-error.
func (s basicService) Get(ctx context.Context, ID string) (Book, error) {
	return s.r.GetByID(ID)
}

// List available items based on limit and offset.
// order takes string in the format "name asc" or "name desc"
// or in combination of multiple fields like "name asc, isbn desc"
// List return all the books in the system
func (s basicService) List(ctx context.Context, order string, limit, offset int) ([]Book, int, error) {
	return s.r.List(order, limit, offset)
}

// Middleware is a service middleware that takes service return service
type Middleware func(Service) Service
