package repository

import "github.com/ishanshre/Book-Review-Platform/internals/models"

// DatabaseRepo consist of all the method available to us to use for database operations
type DatabaseRepo interface {
	AllUsers(limit, offset int) ([]*models.User, error)
	AllReaders(limit, offset int) ([]*models.User, error)

	GetUserByID(id int) (*models.User, error)
	GetGlobalUserByID(id int) (*models.User, error)

	DeleteUser(id int) error
	UpdateUser(u *models.User) error
}
