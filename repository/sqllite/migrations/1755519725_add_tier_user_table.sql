-- +migrate Up
CREATE TABLE user_tiers (
    user_id VARCHAR(32) NOT NULL,
    tier VARCHAR(64) NOT NULL,
    PRIMARY KEY (user_id, tier),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tier) REFERENCES tiers(tier) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE `user_tiers`;