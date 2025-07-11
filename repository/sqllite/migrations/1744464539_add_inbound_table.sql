-- +migrate Up
CREATE TABLE `inbounds` (
        `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `protocol` varchar(32),
        `is_active` boolean default false,
        `domain` varchar(64),
        `vpn_type` varchar(32),
        `port` varchar(8),
        `user_id` varchar(8),
        `tag` varchar(32),
        `is_block` boolean default false,
        `start` TIMESTAMP NOT NULL,
        `end` TIMESTAMP NOT NULL,
        `is_notified` boolean default false,
        `is_assigned` boolean default false,
        `charge_count` INTEGER DEFAULT 0,
        `traffic_usage` INTEGER DEFAULT 0,
        `traffic_limit` INTEGER NOT NULL,
        `country` VARCHAR(128) NOT NULL,
        `is_port_open` BOOLEAN DEFAULT false,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `inbounds`;


