CREATE TABLE journals (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    entry TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);