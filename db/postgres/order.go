package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/db"
	"github.com/kavirajk/bookshop/order"
	_ "github.com/lib/pq"
)

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(driver, source string) (order.Repo, error) {
	db, err := gorm.Open(driver, source)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&order.Order{})
	return &orderRepo{db: db}, nil
}

func (r *orderRepo) get(where ...interface{}) (order.Order, error) {
	var b order.Order
	d := r.db.New()

	if err := d.First(&b, where...).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return order.Order{}, db.ErrNotFound
		}
		return order.Order{}, err
	}
	return b, nil
}

func (r *orderRepo) filter(where ...interface{}) ([]order.Order, error) {
	orders := make([]order.Order, 0)
	d := r.db.New()

	err := d.Find(&orders, where...).Error
	return orders, err
}

func (r *orderRepo) GetByID(ID string) (order.Order, error) {
	return r.get("id=?", ID)
}

func (r *orderRepo) ListByUser(userID string) ([]order.Order, error) {
	return r.filter("created_by_id=?", userID)
}

func (r *orderRepo) Create(u *order.Order) error {
	d := r.db.New()

	if u.ID == "" {
		u.ID = NewID()
	}

	if err := d.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *orderRepo) Save(u *order.Order) error {
	d := r.db.New()

	if err := d.Save(u).Error; err != nil {
		return err
	}
	return nil
}

func (r *orderRepo) Drop() error {
	return r.db.Exec("DELETE FROM ORDERS").Error
}
