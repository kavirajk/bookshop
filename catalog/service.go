package catalog

type Service interface {
	List(tags []string, order string, pageNum, pageSize int) ([]Album, error)
	Search(tag string) ([]Album, error)
	Get(id string) (Album, error)
	Count(tags []string) (int, error)
}
