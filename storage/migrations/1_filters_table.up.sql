CREATE TABLE IF NOT EXISTS filters (
    id UUID UNIQUE NOT NULL,
    type VARCHAR(20) NOT NULL,
    configuration JSON NOT NULL,
    reversed BOOL NOT NULL,
    user_id UUID NOT NULL
);