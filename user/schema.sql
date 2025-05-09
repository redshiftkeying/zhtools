CREATE TABLE users
(
    id         INTEGER PRIMARY KEY,
    name       text NOT NULL,
    password   text,
    phone      text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);