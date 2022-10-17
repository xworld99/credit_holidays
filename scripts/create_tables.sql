CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    balance money NOT NULL default 0
);

CREATE TABLE IF NOT EXISTS services (
    id serial PRIMARY KEY,
    name text NOT NULL,
    desctiption text NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id serial PRIMARY KEY,
    created_at timestamp DEFAULT current_timestamp,
    proofed_at timestamp DEFAULT null,
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    service_id INTEGER REFERENCES services(id) ON DELETE CASCADE,
    amount money NOT NULL,
    is_proofed boolean DEFAULT false
)