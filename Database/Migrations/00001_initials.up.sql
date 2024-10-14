BEGIN;

CREATE TABLE IF NOT EXISTS Users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    role role_type NOT NULL,
    address TEXT NOT NULL REFERENCES Users(id),
    latitude DOUBLE PRECISION NOT NULL REFERENCES Users(id),
    longitude DOUBLE PRECISION NOT NULL REFERENCES Users(id),
    created_by uuid REFERENCES Users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT Now(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_user ON users (email) WHERE archived_at IS NULL;

CREATE TYPE role_type AS ENUM (
    'Admin'
    'Sub-admin'
    'User'
);

CREATE TABLE IF NOT EXISTS Restaurants(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL REFERENCES Users(id),
    longitude DOUBLE PRECISION NOT NULL REFERENCES Users(id),
    created_by uuid REFERENCES Users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT Now(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_restaurant ON restaurants (name, address) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS dishes
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    price INTEGER NOT NULL,
    restaurant_id UUID REFERENCES restaurants (id) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS User_sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES Users(id) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_dish ON dishes (restaurant_id, name) WHERE archived_at IS NULL;


COMMIT;