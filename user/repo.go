package user

type Repo interface {
	GetByID(id string) (User, error)
	GetByUserName(username string) (User, error)
	GetByEmail(email string) (User, error)
	List() ([]User, error)
	Create(user *User) error
	Save(user *User) error
}
