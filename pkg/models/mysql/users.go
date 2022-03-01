package pgql

import (
	"database/sql"

	"github.com/mopeps/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// We'll use the insert method to add a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// We'll use the Authenticate method to verify wether a user exist with
// the provided email address and password. This will return the revelant
// user ID if they do.

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// We'll use the Get method to fetc details for a specific user based
// on their user Id.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
