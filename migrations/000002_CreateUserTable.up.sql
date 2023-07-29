-- Create user table
CREATE TABLE "users" (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    gender gender_enum NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(10),
    profile_pic VARCHAR(255),
    dob DATE,
    document_type document_enum DEFAULT 'Citizenship',
    document_number VARCHAR(20),
    document_front VARCHAR(255),
    document_back VARCHAR(255),
    access_level INTEGER DEFAULT 3,
    is_validated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    last_login TIMESTAMPTZ 
);