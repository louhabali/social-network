-- +migrate Up
CREATE TABLE IF NOT EXISTS tokens (
    user_id INTEGER,
    token TEXT ,
    created_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);