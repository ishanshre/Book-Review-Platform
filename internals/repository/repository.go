package repository

import "github.com/ishanshre/Book-Review-Platform/internals/models"

// DatabaseRepo consist of all the method available to us to use for database operations
type DatabaseRepo interface {
	// User/admin interfaces
	AllUsers(limit, offset int) ([]*models.User, error)
	AllReaders(limit, offset int) ([]*models.User, error)

	GetUserByID(id int) (*models.User, error)
	GetGlobalUserByID(id int) (*models.User, error)

	DeleteUser(id int) error
	UpdateUser(u *models.User) error
	UpdateProfilePic(path string, id int) error

	UpdateLastLogin(id int) error
	Authenticate(username, testPassword string) (int, int, error)
	InsertUser(*models.User) error
	AdminInsertUser(*models.User) error

	GetProfilePersonal(id int) (*models.User, error)

	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)

	ChangePassword(password, email string) error

	// Genre interface
	AllGenre() ([]*models.Genre, error)
	InsertGenre(u *models.Genre) error
	UpdateGenre(u *models.Genre) error
	DeleteGenre(id int) error
	GetGenreByID(id int) (*models.Genre, error)
	GenreExists(title string) (bool, error)

	// Publisher interface
	AllPublishers() ([]*models.Publisher, error)
	InsertPublisher(u *models.Publisher) error
	UpdatePublisher(u *models.Publisher) error
	DeletePublisher(id int) error
	GetPublisherByID(id int) (*models.Publisher, error)
	PublisherExists(name string) (bool, error)

	// Author interface
	AllAuthor() ([]*models.Author, error)
	InsertAuthor(u *models.Author) error
	UpdateAuthor(u *models.Author) error
	DeleteAuthor(id int) error
	GetAuthorByID(id int) (*models.Author, error)

	// Language interface
	AllLanguage() ([]*models.Language, error)
	InsertLanguage(u *models.Language) error
	UpdateLanguage(u *models.Language) error
	DeleteLanguage(id int) error
	GetLanguageByID(id int) (*models.Language, error)
	LanguageExists(language string) (bool, error)

	// book interface
	AllBook() ([]*models.Book, error)
	DeleteBook(id int) error
	InsertBook(u *models.Book) error
	GetBookByID(id int) (*models.Book, error)
	BookIsbnExists(isbn int64) (bool, error)
	UpdateBook(u *models.Book) error
}
