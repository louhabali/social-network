-- +migrate Up
CREATE TABLE IF NOT EXISTS tokens (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    token TEXT UNIQUE,
    created_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);