CREATE TABLE "read_lists" (
    user_id INTEGER,
    book_id INTEGER,
    created_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, book_id),
    CONSTRAINT fk_readLists_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_readLists_book_id FOREIGN KEY (book_id) REFERENCES books(id)
);