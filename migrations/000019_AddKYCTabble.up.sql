CREATE TABLE "kycs" (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    gender gender_enum,
    address VARCHAR(255),
    phone VARCHAR(10),
    profile_pic VARCHAR(255),
    dob DATE,
    document_type document_enum DEFAULT 'Citizenship',
    document_number VARCHAR(50),
    document_front VARCHAR(255),
    document_back VARCHAR(255),
    is_validated BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ,
    CONSTRAINT oneToone_kyc_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);