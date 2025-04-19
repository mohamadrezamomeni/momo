
-- +migrate Up
CREATE TABLE `users` (
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
        `username` varchar(32) UNIQUE NOT NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `lastName` varchar(32),
        `firstName` varchar(32)
);

-- +migrate Down
DROP TABLE `users`;


