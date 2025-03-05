CREATE TABLE token (
  id UUID PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  used_at TIMESTAMP
);

CREATE INDEX idx_token_user_id ON token(user_id);
