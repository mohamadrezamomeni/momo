-- +migrate Up
CREATE TABLE vpn_source (
    country VARCHAR(64) NOT NULL PRIMARY KEY,
    english  VARCHAR(128) NOT NULL
);

-- +migrate Down
DROP TABLE `vpn_source`;