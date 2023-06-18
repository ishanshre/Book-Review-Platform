CREATE TABLE "book_languages" (
    book_id INTEGER NOT NULL,
    language_id INTEGER NOT NULL,
    PRIMARY KEY (book_id, language_id),
    CONSTRAINT fk_book_languages_book_id FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    CONSTRAINT fk_book_languages_language_id FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE
);