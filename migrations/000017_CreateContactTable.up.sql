CREATE TABLE "contacts" (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL,
    phone VARCHAR(15),
    subject VARCHAR(100) NOT NULL,
    message VARCHAR(10000) NOT NULL,
    submitted_at TIMESTAMPTZ,
    ip_address VARCHAR(45),
    browser_info VARCHAR(255),
    referring_page VARCHAR(255)
);