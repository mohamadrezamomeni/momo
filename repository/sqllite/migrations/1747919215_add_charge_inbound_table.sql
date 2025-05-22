-- +migrate Up
CREATE TABLE charges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    status VARCHAR(64) NOT NULL,
    detail TEXT,
    admin_comment TEXT,
    inbound_id INTEGER,
    user_id VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inbound_id) REFERENCES inbounds(inbound_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- +migrate Down
DROP TABLE `charges`;