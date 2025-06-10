-- +migrate Up
CREATE TABLE vpn_source (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(64) NOT NULL,
    english  VARCHAR(128) NOT NULL
);

-- +migrate Down
DROP TABLE `vpn_source`;