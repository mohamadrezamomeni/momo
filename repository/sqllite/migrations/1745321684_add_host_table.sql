-- +migrate Up
CREATE TABLE `host` (
        `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `rank` INTEGER NOT NULL,
        `domain` BOOLEAN default true,
        `port` varchar(8) NOT NULL,
        `status` varchar(16) NOT NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- -migrate Down
DROP TABLE `host`