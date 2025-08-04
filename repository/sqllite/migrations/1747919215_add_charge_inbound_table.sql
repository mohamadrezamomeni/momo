-- +migrate Up
CREATE TABLE charges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    status VARCHAR(64) NOT NULL,
    detail TEXT,
    admin_comment TEXT,
    inbound_id INTEGER,
    user_id VARCHAR(32) NOT NULL,
    package_id VARCHAR(32) NOT NULL,
    country varchar(128) NOT NULL,
    vpn_type varchar(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inbound_id) REFERENCES inbounds(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
    FOREIGN KEY (package_id) REFERENCES vpn_package(id)
);

-- +migrate Down
DROP TABLE `charges`;