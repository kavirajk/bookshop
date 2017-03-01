package postgres_test

import (
	"log"
	"os"
	"testing"

	"github.com/kavirajk/bookshop/db"
	"github.com/kavirajk/bookshop/db/postgres"
	"github.com/kavirajk/bookshop/user"
)

var dbSource string

func init() {
	dbSource = os.Getenv("POSTGRES_TEST_DB_DATASOURCE")
	if dbSource == "" {
		log.Fatal("missing POSTGRES_TEST_DB_DATASOURCE env variable")
	}
}

func setup(t *testing.T) user.Repo {
	repo, err := postgres.NewUserRepo(dbSource)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return repo
}

func TestCreate(t *testing.T) {
	repo := setup(t)
	defer repo.Drop()
	u := user.User{}
	err := repo.Create(&u)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if u.ID == "" {
		t.Errorf("expected non-empty, got empty ID")
	}

	// trying to create with the already existing
	err = repo.Create(&u)
	if err == nil {
		t.Errorf("expected non-empty, got nil error")
	}
}

func TestSave(t *testing.T) {
	repo := setup(t)
	defer repo.Drop()

	t.Run("create via save", func(t *testing.T) {
		u := user.User{}
		err := repo.Save(&u)
		if err == nil {
			t.Errorf("exptected non-nil error, got nil")
		}

	})
}

func TestGetByID(t *testing.T) {
	repo := setup(t)
	defer repo.Drop()
	t.Run("invalid", func(t *testing.T) {
		_, err := repo.GetByID("1")
		if err != db.ErrNotFound {
			t.Errorf("expected NotFound, got %v", err)
		}
	})

	var u user.User
	repo.Create(&u) // no-need error check. Verified in create test
	t.Run("valid", func(t *testing.T) {
		_, err := repo.GetByID(u.ID)
		if err != nil {
			t.Errorf("expected non-nil error, got %v", err)
		}
	})
}

func TestGetByEmail(t *testing.T) {
	repo := setup(t)
	defer repo.Drop()
	t.Run("invalid", func(t *testing.T) {
		_, err := repo.GetByEmail("not-exists@golang.org")
		if err != db.ErrNotFound {
			t.Errorf("expected NotFound, got %v", err)
		}
	})

	u := user.User{Email: "chandler@golang.org"}
	repo.Create(&u) // no-need error check. Verified in create test
	t.Run("valid", func(t *testing.T) {
		_, err := repo.GetByEmail("chandler@golang.org")
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})
}

func TestGetByToken(t *testing.T) {
	repo := setup(t)
	defer repo.Drop()
	t.Run("invalid", func(t *testing.T) {
		_, err := repo.GetByToken("xxxxx")
		if err != db.ErrNotFound {
			t.Errorf("expected NotFound, got %v", err)
		}
	})
	token := "xghfghfghfgh"
	u := user.User{Email: "chandler@golang.org", AuthToken: token}
	repo.Create(&u) // no-need error check. Verified in create test
	t.Run("valid", func(t *testing.T) {
		_, err := repo.GetByToken(token)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})
}
