ALTER TABLE "reviews"
ADD CONSTRAINT uc_review_user_book UNIQUE (user_id, book_id);