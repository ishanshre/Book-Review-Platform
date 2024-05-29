CREATE TABLE "publishers" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(10000),
    pic VARCHAR(255),
    address VARCHAR(255),
    phone VARCHAR(10),
    email VARCHAR(255),
    website VARCHAR(255),
    established_date INTEGER,
    latitude VARCHAR(20),
    longitude VARCHAR(20)
)