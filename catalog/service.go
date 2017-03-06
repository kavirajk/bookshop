package catalog

import "context"

type Service interface {
	// Search books based on free text
	Search(ctx context.Context, query string) ([]Book, error)

	// Get details about single book
	Get(ctx context.Context, id string) (Book, error)
}

type basicService struct {
	r Repo
}

// NewCatalogService return basic Service implementation.
func NewCatalogService(r Repo) Service {
	return basicService{r}
}

// Search return books that matches with query.
func (s basicService) Search(ctx context.Context, query string) ([]Book, error) {
	return nil, nil
}

// Get return a book for the matched ID. Empty book incase of non-error.
func (s basicService) Get(ctx context.Context, ID string) (Book, error) {
	return Book{}, nil
}

type Middleware func(Service) Service
