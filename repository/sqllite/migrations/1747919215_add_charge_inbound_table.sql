-- +migrate Up
CREATE TABLE charges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    status VARCHAR(64) NOT NULL,
    detail TEXT,
    admin_comment TEXT,
    inbound_id INTEGER, 
    FOREIGN KEY (inbound_id) REFERENCES inbounds(inbound_id)
);

-- +migrate Down
DROP TABLE `charges`;