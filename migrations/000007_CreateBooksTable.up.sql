CREATE TABLE "books" (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(10000),
    cover VARCHAR(255),
    isbn BIGINT UNIQUE NOT NULL CHECK (isbn >= 1000000000000 AND isbn <= 9999999999999),
    published_date TIMESTAMPTZ,
    paperback INTEGER,
    is_active BOOLEAN DEFAULT TRUE,
    added_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    publisher_id int
);