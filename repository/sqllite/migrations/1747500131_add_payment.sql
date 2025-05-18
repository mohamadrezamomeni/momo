-- +migrate Up
CREATE TABLE `payment` (
       `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `is_settled` BOOLEAN default false,
        `user_id` varchar(32) NOT NULL,
        `start` TIMESTAMP NOT NULL,
        `end` TIMESTAMP NOT NULL,
        `inbound_id` INTEGER NOT NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `payment`;