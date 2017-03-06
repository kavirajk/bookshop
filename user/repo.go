package user

// Repo abstracts all the persistant storage operations of User service.
type Repo interface {
	Create(user *User) error
	Save(user *User) error
	GetByID(id string) (User, error)
	GetByUserName(username string) (User, error)
	GetByEmail(email string) (User, error)
	GetByToken(token string) (User, error)
	GetByResetKey(email string) (User, error)
	List(order string, limit, offset int) (users []User, total int, err error)
	Drop() error
}
