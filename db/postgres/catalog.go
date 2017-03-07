package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/catalog"
	"github.com/kavirajk/bookshop/db"
	_ "github.com/lib/pq"
)

type catalogRepo struct {
	db *gorm.DB
}

func NewCatalogRepo(source string) (catalog.Repo, error) {
	db, err := gorm.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&catalog.Book{})
	return &catalogRepo{db: db}, nil
}

func (r *catalogRepo) get(where ...interface{}) (catalog.Book, error) {
	var b catalog.Book
	d := r.db.New()

	if err := d.First(&b, where...).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return catalog.Book{}, db.ErrNotFound
		}
		return catalog.Book{}, err
	}
	return b, nil
}

func (r *catalogRepo) filter(where ...interface{}) ([]catalog.Book, error) {
	books := make([]catalog.Book, 0)
	d := r.db.New()

	err := d.Find(&books, where...).Error
	return books, err
}

func (r *catalogRepo) GetByID(ID string) (catalog.Book, error) {
	return r.get("id=?", ID)
}

func (r *catalogRepo) GetByISBN(ISBN string) (catalog.Book, error) {
	return r.get("isbn=?", ISBN)
}

func (r *catalogRepo) ListByAuthor(authorID string) ([]catalog.Book, error) {
	// TODO(kaviraj)
	// Query from intermediate table
	return nil, nil
}

func (r *catalogRepo) GetByToken(token string) (catalog.Book, error) {
	return r.get("auth_token=?", token)
}

func (r *catalogRepo) GetByResetKey(key string) (catalog.Book, error) {
	return r.get("reset_key=?", key)
}

func (r *catalogRepo) List(order string, limit, offset int) ([]catalog.Book, int, error) {
	catalogs := make([]catalog.Book, 0)
	db := r.db.New()

	var total int
	if err := db.Model(&catalog.Book{}).Order(order).Count(&total).Error; err != nil {
		return catalogs, 0, err
	}

	err := db.Order(order).Limit(limit).Offset(offset).Find(&catalogs).Error
	return catalogs, total, err
}

func (r *catalogRepo) Create(u *catalog.Book) error {
	d := r.db.New()

	if u.ID == "" {
		u.ID = NewID()
	}

	if err := d.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *catalogRepo) Save(u *catalog.Book) error {
	d := r.db.New()

	if err := d.Save(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *catalogRepo) Drop() error {
	return r.db.Exec("DELETE FROM CATALOGS").Error
}
