
CREATE TABLE statuses (
  id SERIAL PRIMARY KEY,
  status VARCHAR(20) NOT NULL
);
INSERT INTO statuses (status) VALUES ('NEW'), ('PROCESSING'), ('INVALID'), ('PROCESSED'); 
