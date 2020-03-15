CREATE TABLE IF NOT EXISTS filters (
    name VARCHAR(30) UNIQUE NOT NULL,
    type VARCHAR(20) NOT NULL,
    configuration JSON NOT NULL,
    reversed BOOL NOT NULL,
    user_id VARCHAR(20) NOT NULL
)