CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    UUID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    FirstName VARCHAR(255),
    LastName VARCHAR(255),
    Age INT,
    Address VARCHAR(255),
    Email VARCHAR(255)
);