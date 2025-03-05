CREATE TABLE accrual_statuses (
  id SERIAL PRIMARY KEY,
  status VARCHAR(20) NOT NULL
);

INSERT INTO accrual_statuses (status) VALUES ('REGISTERED'), ('INVALID'), ('PROCESSING'), ('PROCESSED'); 
