CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL, -- Store hashed password
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
    description TEXT
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
    is_replacement_needed BOOLEAN DEFAULT FALSE
);

-- Investment Tracking Table
CREATE TABLE item_investments (
    id SERIAL PRIMARY KEY,
    item_id INTEGER REFERENCES items(id),
    initial_price DECIMAL(10, 2),
    current_value DECIMAL(10, 2),
    depreciation_rate DECIMAL(5, 2),
    last_depreciation_date DATE
);
