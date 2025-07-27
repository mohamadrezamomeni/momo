-- +migrate Up
CREATE TABLE `vpns` (
        `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `domain` varchar(64) NOT NULL,
        `is_active` boolean default true,
        `api_port` varchar(8) NOT NULL,
        `vpn_type` varchar(32) NOT NULL,
        `user_count` int NOT NULL,
        `country` varchar(128) NOT NULL,
        `start_port` varchar(64) NOT NULL,
        `end_port` varchar(64) NOT NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `vpns`;


