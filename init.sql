CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    vk_id INTEGER,
    first_name varchar(255),
    last_name varchar(255)
);
CREATE TABLE IF NOT EXISTS types (
    id SERIAL PRIMARY KEY,
    title varchar(255)
);
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER references users(id),
    title varchar(255),
    description varchar(255),
    type_id INTEGER references types(id),
    location geography(POINT),
    radius INTEGER,
    status BOOLEAN,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);
CREATE TABLE IF NOT EXISTS points (
    id SERIAL PRIMARY KEY,
    user_id INTEGER references users(id),
    title varchar(255),
    location geography(POINT),
    radius INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);