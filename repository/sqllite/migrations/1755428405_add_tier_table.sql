-- +migrate Up
CREATE TABLE tiers (
    `name` VARCHAR(64) NOT NULL PRIMARY KEY,
    `is_default` BOOLEAN DEFAULT false
);

-- +migrate Down
DROP TABLE `tiers`;