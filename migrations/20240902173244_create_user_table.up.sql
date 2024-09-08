-- migration up for create_user_table
CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    username VARCHAR(40) UNIQUE NOT NULL CHECK (username <> ''),
    email VARCHAR(40) UNIQUE NOT NULL CHECK (email <> ''),
    password VARCHAR(40),
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp
);