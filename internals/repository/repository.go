package repository

import "github.com/ishanshre/Book-Review-Platform/internals/models"

// DatabaseRepo consist of all the method available to us to use for database operations
type DatabaseRepo interface {
	// User/admin interfaces
	AllUsers(limit, offset int) ([]*models.User, error)
	AllReaders(limit, offset int) ([]*models.User, error)

	GetUserByID(id int) (*models.User, error)
	GetGlobalUserByID(id int) (*models.User, error)
	GetGlobalUserByIDAny(id int) (*models.User, error)

	DeleteUser(id int) error
	UpdateUser(u *models.User) error
	UpdateProfilePic(path string, id int) error

	UpdateLastLogin(id int) error
	Authenticate(username, testPassword string) (int, int, bool, error)
	InsertUser(*models.User) error
	AdminInsertUser(*models.User) error

	GetProfilePersonal(id int) (*models.User, error)

	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)

	ChangePassword(password, email string) error
	UserListFilter(limit, page int, searchKey, sort string) (*models.AdminUserListApi, error)

	// User Kyc
	GetKycByUserID(user_id int) (*models.Kyc, error)
	GetUserWithKyc(id int) (*models.UserKycData, error)
	UpdateDocument(front_path, back_path string, id int) error
	AdminKycUpdate(update *models.Kyc) error
	PublicKycUpdate(update *models.Kyc) error

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
	PublisherExistsID(id int) (bool, error)
	GetPublisherWithBookByID(publisher_id int) (*models.PublisherWithBooksData, error)
	AllPublishersFilter(limit, page int, searchKey, sort string) (*models.AdminPublisherListApi, error)

	// Author interface
	AllAuthor() ([]*models.Author, error)
	InsertAuthor(u *models.Author) error
	UpdateAuthor(u *models.Author) error
	DeleteAuthor(id int) error
	GetAuthorByID(id int) (*models.Author, error)
	GetAuthorFullNameByID(id int) (*models.Author, error)
	TotalAuthors() (int, error)
	AllAuthorsFilter(limit, page int, search, order string) (*models.AuthorApiFilter, error)
	GetAuthorWithBooks(id int) (*models.AuthorBookData, error)

	// Language interface
	AllLanguage() ([]*models.Language, error)
	InsertLanguage(u *models.Language) error
	UpdateLanguage(u *models.Language) error
	DeleteLanguage(id int) error
	GetLanguageByID(id int) (*models.Language, error)
	LanguageExists(language string) (bool, error)

	// book interface
	AllBook() ([]*models.Book, error)
	AllBookData(limit, page int) ([]*models.Book, error)
	AllBookDataRandom() ([]*models.Book, error)
	AllBookRandomPage(limit, page int) ([]*models.Book, error)
	DeleteBook(id int) error
	InsertBook(u *models.Book) error
	GetBookByID(id int) (*models.Book, error)
	GetBookByISBN(isbn int64) (*models.Book, error)
	BookIsbnExists(isbn int64) (bool, error)
	UpdateBook(u *models.Book) error
	GetBookTitleByID(id int) (*models.Book, error)
	TotalBooks() (int, error)
	AllBooksFilter(limit, page int, searchKey, sort string) (*models.BookApiFilter, error)
	BookDetailWithAuthorPublisherWithIsbn(isbn int64) (*models.BookInfoData, error)
	AllRecentBooks(limit, page int) ([]*models.Book, error)

	CalculateLastPage(limit, total int) int

	// book author interface
	AllBookAuthor() ([]*models.BookAuthor, error)
	DeleteBookAuthor(book_id, author_id int) error
	GetBookAuthorByID(book_id, author_id int) (*models.BookAuthor, error)
	GetBookAuthorByBookID(book_id int) ([]*models.BookAuthor, error)
	BookAuthorExists(book_id, author_id int) (bool, error)
	UpdateBookAuthor(u *models.BookAuthor, book_id, author_id int) error
	InsertBookAuthor(u *models.BookAuthor) error
	BookAuthorListFilter(limit, page int, searchKey, sort string) (*models.BookAuthorListApi, error)

	// book genre interface
	AllBookGenre() ([]*models.BookGenre, error)
	DeleteBookGenre(book_id, genre_id int) error
	GetBookGenreByID(book_id, genre_id int) (*models.BookGenre, error)
	BookGenreExists(book_id, genre_id int) (bool, error)
	UpdateBookGenre(u *models.BookGenre, book_id, genre_id int) error
	InsertBookGenre(u *models.BookGenre) error
	GetGenresFromBookID(book_id int) ([]*models.Genre, error)
	GetAllBooksByGenre(limit, page int, searchKey, sort, genre string) (*models.BookApiFilter, error)

	// Book Language interface
	AllBookLanguage() ([]*models.BookLanguage, error)
	DeleteBookLanguage(book_id, language_id int) error
	GetBookLanguageByID(book_id, language_id int) (*models.BookLanguage, error)
	BookLanguageExists(book_id, language_id int) (bool, error)
	UpdateBookLanguage(u *models.BookLanguage, book_id, language_id int) error
	InsertBookLanguage(u *models.BookLanguage) error
	GetLanguagesFromBookID(book_id int) ([]*models.Language, error)
	GetAllBooksByLanguage(limit, page int, searchKey, sort, language string) (*models.BookApiFilter, error)

	// ReadList interface
	AllReadList() ([]*models.ReadList, error)
	ReadListExists(user_id, book_id int) (bool, error)
	InsertReadList(u *models.ReadList) error
	GetReadListByID(user_id, book_id int) (*models.ReadList, error)
	DeleteReadList(user_id, book_id int) error
	UpdateReadList(u *models.ReadList, book_id, user_id int) error
	ReadListCount(user_id int) (int, error)
	GetAllBooksFromReadListByUserId(limit, page, user_id int, searchKey, sort string) (*models.BookApiFilter, error)
	ReadListFilter(limit, page int, searchKey, sort string) (*models.ReadListFilterApi, error)

	// BuyList interface
	AllBuyList() ([]*models.BuyList, error)
	BuyListExists(user_id, book_id int) (bool, error)
	InsertBuyList(u *models.BuyList) error
	GetBuyListByID(user_id, book_id int) (*models.BuyList, error)
	DeleteBuyList(user_id, book_id int) error
	UpdateBuyList(u *models.BuyList, book_id, user_id int) error
	BuyListCount(user_id int) (int, error)
	GetAllBooksFromBuyListByUserId(limit, page, user_id int, searchKey, sort string) (*models.BookApiFilter, error)
	BuyListFilter(limit, page int, searchKey, sort string) (*models.BuyListFilterApi, error)

	// Follower Interface
	AllFollowers() ([]*models.Follower, error)
	FollowerExists(u *models.Follower) (bool, error)
	InsertFollower(u *models.Follower) error
	GetFollowerByID(user_id, author_id int) (*models.Follower, error)
	DeleteFollower(user_id, author_id int) error
	UpdateFollower(u *models.Follower, user_id, author_id int) error
	FollowerCount(user_id int) (int, error)
	GetAllFollowingsByUserId(user_id int) ([]*models.Author, error)
	FollowerFilter(limit, page int, searchKey, sort string) (*models.FollowerFilterApi, error)

	// Review interface
	AllReviews() ([]*models.Review, error)
	ReviewExists(u *models.Review) (bool, error)
	InsertReview(u *models.Review) error
	GetReviewByID(id int) (*models.Review, error)
	GetReviewByUserID(id int) (*models.Review, error)
	DeleteReview(id int) error
	UpdateReview(u *models.Review) error
	GetReviewsByBookID(bookID int) ([]*models.Review, error)
	UpdateReviewBook(update *models.Review) error

	// Contact interface
	AllContacts() ([]*models.Contact, error)
	GetContactByID(id int) (*models.Contact, error)
	DeleteContact(id int) error
	InsertContact(*models.Contact) error

	// request_books interface
	InsertRequestedBook(i *models.RequestedBook) error
	AllRequestBooks() ([]*models.RequestedBook, error)
	DeleteRequestBooks(id int) error
	GetRequestBookById(id int) (*models.RequestedBook, error)
}
