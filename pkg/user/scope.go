package user

import "github.com/jinzhu/gorm"

// Scope represents gorm scope to chain the db query.
type Scope func(*gorm.DB) *gorm.DB
