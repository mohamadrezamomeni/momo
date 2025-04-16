
-- +migrate Up
CREATE TABLE `users` (
        `id` varchar(255)  PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
        `email` varchar(32) UNIQUE NOT NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `lastName` varchar(32),
        `firstName` varchar(32)
);

-- +migrate Down
DROP TABLE `users`;


