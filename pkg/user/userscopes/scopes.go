package userscopes

import "github.com/jinzhu/gorm"

// GetByID is user scope to get user by ID.
func GetByID(ID int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id=?", ID)
	}
}
