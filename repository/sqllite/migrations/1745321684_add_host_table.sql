-- +migrate Up
CREATE TABLE `hosts` (
       `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `domain` BOOLEAN default true,
        `port` varchar(8) NOT NULL,
        `status` varchar(16) NOT NULL,
        `start_range_port` int(8) NOT NULL,
        `end_range_port` int(8) NOT NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `hosts`;