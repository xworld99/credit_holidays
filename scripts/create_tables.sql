CREATE TABLE IF NOT EXISTS users (
    -- unique id, not serial because of task text
    id int PRIMARY KEY,
    -- current active balance of user
    balance money NOT NULL default 0.0 CHECK(balance >= 0.0),
    -- reserved balance of user
    frozen_balance money DEFAULT 0.0 CHECK(frozen_balance >= 0.0)
);

CREATE TABLE IF NOT EXISTS services (
    -- unique id
    id serial PRIMARY KEY,
    -- name of service
    name text NOT NULL UNIQUE,
    -- should this service be proofed before changing balance of user
    confirmation_needed boolean default false,
    -- text description
    desctiption text NOT NULL,
    -- specify if balance should be increased or decreased
    service_type text NOT NULL CHECK(service_type in ('accrual', 'withdraw')) DEFAULT 'withdraw'
);

INSERT INTO services(name, description, confirmation_needed, service_type)
VALUES (('deposit', 'replenishment of the balance from an external source', false, 'accrual'),
        ('withdraw', 'withdrawal of money to an external source', true, 'withdraw'),
        ('transfer', 'transfer money from one user to another one', true, 'withdraw'),
        ('subscription', 'withdrawal of money to an external source', true, 'withdraw'),
        ('payment', 'user buy something using his virtual balance', true, 'withdraw')
        ('cashback', 'cashback on user balance', false, 'accrual'));

CREATE TABLE IF NOT EXISTS orders (
    -- unique id
    id serial PRIMARY KEY,
    -- date of orders initiation
    created_at timestamp DEFAULT current_timestamp,
    -- date of orders proof
    proofed_at timestamp DEFAULT null,
    user_id integer REFERENCES users(id) ON DELETE CASCADE,
    service_id INTEGER REFERENCES services(id) ON DELETE CASCADE,
    amount money NOT NULL CHECK(amount >= 0.0),
    status text DEFAULT 'in_progress' CHECK(status in ('in_progress', 'success', 'declined'))
);