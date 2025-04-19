-- +migrate Up
CREATE TABLE `inbounds` (
        `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `protocol` varchar(32),
        `is_available` boolean default false,
        `domain` varchar(64),
        `vpn_type` varchar(32),
        `port` varchar(8),
        `user_id` varchar(8),
        `tag` varchar(32),
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `inbounds`;


