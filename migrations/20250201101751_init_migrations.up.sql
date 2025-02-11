CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(100) NOT NULL,
    password VARCHAR(60) NOT NULL, 
    email VARCHAR(150) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE token (
  id UUID PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  used_at TIMESTAMP
);

CREATE INDEX idx_token_user_id ON token(user_id);

CREATE TABLE statuses (
  id SERIAL PRIMARY KEY,
  status VARCHAR(20) NOT NULL
);

INSERT INTO statuses (status) VALUES ('NEW'), ('PROCESSING'), ('INVALID'), ('PROCESSED'); 

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  status_id INTEGER NOT NULL REFERENCES statuses(id),
  creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  number BIGINT NOT NULL,
  uploaded_at TIMESTAMP
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status_id ON orders(status_id);
