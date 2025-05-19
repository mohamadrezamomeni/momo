-- +migrate Up
CREATE TABLE `vpn_package` (
       `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `price_tilte` varchar(64) NOT NULL,
        `price` INTEGER NOT NULL,
        `days` INTEGER NOT NULL,
        `months` INTEGER NOT NULL,
        `traffic_limit` INTEGER NOT NULL,
        `traffic_limit_title` varchar NOT NULL,
        `is_active`  BOOLEAN DEFAULT true,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `vpn_package`;