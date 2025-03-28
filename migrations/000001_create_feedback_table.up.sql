CREATE TABLE feedback (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(50),
    subject VARCHAR(50),
    message TEXT,
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);