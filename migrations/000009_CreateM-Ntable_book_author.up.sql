CREATE TABLE "book_authors" (
    book_id INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    PRIMARY KEY (book_id, author_id),
    CONSTRAINT fk_books_authors_book_id FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    CONSTRAINT fk_books_authors_author_id FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE CASCADE
);