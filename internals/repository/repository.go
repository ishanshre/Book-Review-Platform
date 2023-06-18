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
	GetAuthorFullNameByID(id int) (*models.Author, error)

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
	GetBookTitleByID(id int) (*models.Book, error)

	// book author interface
	AllBookAuthor() ([]*models.BookAuthor, error)
	DeleteBookAuthor(book_id, author_id int) error
	GetBookAuthorByID(book_id, author_id int) (*models.BookAuthor, error)
	BookAuthorExists(book_id, author_id int) (bool, error)
	UpdateBookAuthor(u *models.BookAuthor, book_id, author_id int) error
	InsertBookAuthor(u *models.BookAuthor) error

	// book genre interface
	AllBookGenre() ([]*models.BookGenre, error)
	DeleteBookGenre(book_id, genre_id int) error
	GetBookGenreByID(book_id, genre_id int) (*models.BookGenre, error)
	BookGenreExists(book_id, genre_id int) (bool, error)
	UpdateBookGenre(u *models.BookGenre, book_id, genre_id int) error
	InsertBookGenre(u *models.BookGenre) error

	// Book Language interface
	AllBookLanguage() ([]*models.BookLanguage, error)
	DeleteBookLanguage(book_id, language_id int) error
	GetBookLanguageByID(book_id, language_id int) (*models.BookLanguage, error)
	BookLanguageExists(book_id, language_id int) (bool, error)
	UpdateBookLanguage(u *models.BookLanguage, book_id, language_id int) error
	InsertBookLanguage(u *models.BookLanguage) error
}
