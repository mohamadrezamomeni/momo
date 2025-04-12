-- +migrate Up
CREATE TABLE `inbounds` (
        `id` varchar(32) PRIMARY KEY DEFAULT (
        lower(
            printf('%08x-%04x-%04x-%04x-%012x',
                abs(random()) % 4294967296,           
                abs(random()) % 65536,                
                0x4000 | (abs(random()) % 4096),      
                0x8000 | (abs(random()) % 16384),     
                abs(random()) % 281474976710656       
            )
        )
        ),
        `protocol` varchar(32),
        `isAvailable` boolean default true,
        `domain` varchar(64),
        `port` varchar(8),
        `user_id` varchar(8),
        `tag` varchar(32),
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `inbounds`;


