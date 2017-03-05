package catalog

import "context"

type Service interface {
	List(ctx context.Context, tags []string, order string, limit, offset int) ([]Book, error)
	Search(ctx context.Context, tag string) ([]Book, error)
	Get(ctx context.Context, id string) (Book, error)
	Count(ctx context.Context, tags []string) (int, error)
}
