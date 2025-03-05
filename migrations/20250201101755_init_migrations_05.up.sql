CREATE TABLE balance (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  current_balance DECIMAL NOT NULL,
  withdrawn DECIMAL NULL,
  updated VARCHAR(150) NULL
);