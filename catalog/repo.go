package catalog

// Repo abstracts all the persistant storage operations of Catalog Service
type Repo interface {
	Create(book *Book) error
	Save(book *Book) error
	GetByID(ID string) (Book, error)
	List(order string, limit, offset int) ([]Book, int, error)
	Search(name string) ([]Book, error)
	GetByISBN(ISBN string) (Book, error)
	ListByAuthor(authorID string) ([]Book, error)
	Drop() error
}
