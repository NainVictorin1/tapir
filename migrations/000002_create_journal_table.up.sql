CREATE TABLE journal (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(50),
    subject VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);