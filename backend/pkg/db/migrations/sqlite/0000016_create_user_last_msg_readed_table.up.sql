-- +migrate Up
CREATE TABLE user_lastMsg_read(
    user_id INTEGER NOT NULL,
    chat_id INTEGER NOT NULL,
    last_read_message_id INTEGER,
    PRIMARY KEY (user_id, chat_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    FOREIGN KEY (last_read_message_id) REFERENCES messages(id)
);