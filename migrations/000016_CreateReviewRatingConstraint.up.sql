ALTER TABLE "reviews"
ADD CONSTRAINT review_rating_limit_check
CHECK (
    rating >= 1.0
    AND rating <= 5.0
);