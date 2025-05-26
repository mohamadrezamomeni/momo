-- +migrate Up
CREATE TABLE `events` (
       `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `name` varchar(64) NOT NULL,
        `data` TEXT NOT NULL,
        `is_notification_processed` BOOLEAN DEFAULT false,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `events`;