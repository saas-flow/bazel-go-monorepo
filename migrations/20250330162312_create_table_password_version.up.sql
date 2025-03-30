
CREATE TABLE IF NOT EXISTS password_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    version INT UNIQUE NOT NULL
);