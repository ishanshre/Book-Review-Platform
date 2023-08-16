CREATE TABLE "request_books" (
    id SERIAL PRIMARY KEY,
    book_title VARCHAR(255),
    author VARCHAR(255),
    requested_email VARCHAR(255),
    requested_date TIMESTAMPTZ
) 