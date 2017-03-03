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

func TestList(t *testing.T) {
	repo := setup(t)
	defer repo.Drop()
	users := []user.User{
		{Email: "test1@bookshop.com", Username: "test1"},
		{Email: "test2@bookshop.com", Username: "test2"},
		{Email: "test3@bookshop.com", Username: "test3"},
		{Email: "test4@bookshop.com", Username: "test4"},
	}
	for i := range users {
		if err := repo.Create(&users[i]); err != nil {
			t.Errorf("expected nil error, got %v\n", err)
		}
	}

	t.Run("test limit", func(t *testing.T) {
		us, total, err := repo.List("", 2, 0)
		if err != nil {
			t.Errorf("expected nil error, got %v\n", err)
		}
		if total != 4 {
			t.Errorf("expected total 4, got %v\n", total)
		}
		if len(us) != 2 {
			t.Errorf("expected return length 2, got %v\n", len(us))
		}
	})
	t.Run("test offset", func(t *testing.T) {
		us, total, err := repo.List("", 5, 1) // starting from offset 1

		if err != nil {
			t.Errorf("expected nil error, got %v\n", err)
		}
		if total != 4 {
			t.Errorf("expected total 4, got %v\n", total)
		}
		if len(us) != 3 {
			t.Errorf("expected return length 3, got %v\n", len(us))
		}
	})
	t.Run("test ordering", func(t *testing.T) {
		us, _, err := repo.List("username", 3, 0)

		if err != nil {
			t.Errorf("expected nil error, got %v\n", err)
		}

		if us[0].Username != "test1" {
			t.Errorf("ordering failed. expected test1, got %v\n", us[0].Username)
		}
		us, _, err = repo.List("username desc", 3, 0)
		if us[0].Username != "test4" {
			t.Errorf("ordering failed. expected test1, got %v\n", us[0].Username)
		}

	})

}
