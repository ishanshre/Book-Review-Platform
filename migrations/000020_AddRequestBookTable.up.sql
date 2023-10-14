CREATE TABLE "request_books" (
    id SERIAL PRIMARY KEY,
    book_title VARCHAR(255),
    author VARCHAR(255),
    requested_by INTEGER,
    requested_date TIMESTAMPTZ,
    is_added BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_requested_book_user FOREIGN KEY (requested_by) REFERENCES users(id) ON DELETE CASCADE
) 