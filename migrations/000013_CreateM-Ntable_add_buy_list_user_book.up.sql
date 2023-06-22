CREATE TABLE "buy_lists" (
    user_id INTEGER,
    book_id INTEGER,
    created_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, book_id),
    CONSTRAINT fk_buyLists_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_buyLists_book_id FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);