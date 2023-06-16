ALTER TABLE books
ADD CONSTRAINT fk_book_publisher_table
FOREIGN KEY (publisher_id) REFERENCES publishers(id);