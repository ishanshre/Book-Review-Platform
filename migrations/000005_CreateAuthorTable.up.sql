CREATE TABLE "authors" (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    bio VARCHAR(10000),
    date_of_birth INTEGER,
    email VARCHAR(255),
    country_of_origin VARCHAR(50),
    avatar VARCHAR(255)
);