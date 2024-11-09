CREATE TYPE status_enum AS ENUM (
	'active',
	'deleted'
)

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL, -- Store hashed password
	email VARCHAR UNIQUE NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    session_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Categories Table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
	status status_enum DEFAULT 'active',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Items Table
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    photo_url VARCHAR(255),
    price DECIMAL(10, 2),
    purchase_date DATE,
    total_usage_days INTEGER DEFAULT 0,
    is_replacement_needed BOOLEAN DEFAULT FALSE,
	status status_enum DEFAULT 'active',
	depreciated_rate INTEGER,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Investment Tracking Table
CREATE TABLE item_investments (
    id SERIAL PRIMARY KEY,
    item_id INTEGER REFERENCES items(id),
    initial_price DECIMAL(10, 2),
    current_value DECIMAL(10, 2),
    last_depreciation_date DATE
);

SELECT * FROM users
SELECT * FROM sessions
SELECT * FROM categories
SELECT * FROM items
SELECT * FROM item_investments

SELECT SUM(initial_price) AS total_investment, SUM(current_value) AS depreciated_value FROM item_investments;