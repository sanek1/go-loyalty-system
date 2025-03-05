CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL, 
    email VARCHAR(150) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated VARCHAR(150) NULL
);

CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_email ON users(email);
