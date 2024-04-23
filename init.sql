CREATE TABLE IF NOT EXISTS tax_rates (
  id SERIAL PRIMARY KEY,
  minimum_salary INT NOT NULL,
  maximum_salary INT NULL,
  rate INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tax_rates (minimum_salary, maximum_salary, rate) VALUES 
(0, 150000, 0),
(150001, 500000, 10), --35,000 | 35,000
(500001, 1000000, 15), -- 75,000 | 110,000
(1000001, 2000000, 20), -- 200,000 | 310,000
(2000001, 0, 35);

CREATE TYPE deducation_type AS ENUM ('Personal', 'Donation','K-Receipt');

CREATE TABLE IF NOT EXISTS tax_deductions (
  id SERIAL PRIMARY KEY,
  type deducation_type NOT NULL,
  minimum_amount INT NOT NULL,
  maximum_amount INT NULL,
  amount INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tax_deductions (type, minimum_amount, maximum_amount, amount) VALUES 
('Personal', 10000, 100000, 60000),
('Donation', 0, 100000, 100000),
('K-Receipt', 0, 100000, 50000);