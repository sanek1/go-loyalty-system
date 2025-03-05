
CREATE TABLE accrual (
  id SERIAL PRIMARY KEY,
  order_id INTEGER  NOT NULL REFERENCES orders(id),
  status_id INTEGER NOT NULL REFERENCES accrual_statuses(id),
  accrual DECIMAL  NULL,
  updated VARCHAR(150) NULL
);
