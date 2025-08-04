-- +migrate Up
ALTER TABLE inbounds
ADD COLUMN vpn_id INTEGER REFERENCES vpns(id);

-- +migrate Down
CREATE TABLE `inbounds_new` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `protocol` VARCHAR(32),
    `is_active` BOOLEAN DEFAULT false,
    `domain` VARCHAR(64),
    `vpn_type` VARCHAR(32),
    `port` VARCHAR(8),
    `user_id` VARCHAR(8),
    `tag` VARCHAR(32),
    `is_block` BOOLEAN DEFAULT false,
    `start` TIMESTAMP NOT NULL,
    `end` TIMESTAMP NOT NULL,
    `is_notified` BOOLEAN DEFAULT false,
    `is_assigned` BOOLEAN DEFAULT false,
    `traffic_usage` INTEGER DEFAULT 0,
    `traffic_limit` INTEGER NOT NULL,
    `country` VARCHAR(128) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO inbounds_new (
    id, protocol, is_active, domain, vpn_type, port, user_id, tag,
    is_block, start, end, is_notified, is_assigned,
    traffic_usage, traffic_limit, country, created_at, updated_at
)
SELECT
    id, protocol, is_active, domain, vpn_type, port, user_id, tag,
    is_block, start, end, is_notified, is_assigned,
    traffic_usage, traffic_limit, country, created_at, updated_at
FROM inbounds;

DROP TABLE inbounds;

ALTER TABLE inbounds_new RENAME TO inbounds;