CREATE TABLE "followers" (
    user_id INTEGER,
    author_id INTEGER,
    followed_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, author_id),
    CONSTRAINT fk_followers_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_followers_author_id FOREIGN KEY (author_id) REFERENCES authors(id)
);