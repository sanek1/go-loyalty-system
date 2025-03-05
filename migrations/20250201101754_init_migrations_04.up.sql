CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  status_id INTEGER NOT NULL REFERENCES statuses(id),
  creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  number BIGINT NOT NULL,
  uploaded_at TIMESTAMP,
  updated VARCHAR(150) NULL
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status_id ON orders(status_id);
