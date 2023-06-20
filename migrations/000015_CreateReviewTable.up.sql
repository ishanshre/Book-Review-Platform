CREATE TABLE "reviews" (
    id SERIAL PRIMARY KEY,
    rating NUMERIC(2,1) NOT NULL,
    body VARCHAR(10000),
    book_id INTEGER,
    user_id INTEGER,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_review_book FOREIGN KEY (book_id) REFERENCES books(id), 
    CONSTRAINT fk_review_user FOREIGN KEY (user_id) REFERENCES users(id)
);