package user

// Repo abstracts all the persistant storage operations of User service.
type Repo interface {
	GetByID(id string) (User, error)
	GetByUserName(username string) (User, error)
	GetByEmail(email string) (User, error)
	GetByToken(token string) (User, error)
	GetByResetKey(email string) (User, error)
	List() ([]User, error)
	Create(user *User) error
	Save(user *User) error
}
