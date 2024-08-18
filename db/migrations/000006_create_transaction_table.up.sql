CREATE TABLE IF NOT EXISTS transaction (
  id SERIAL NOT NULL PRIMARY KEY,
  patient_id INT NOT NULL,
  patient_address_id INT NOT NULL,
  status VARCHAR(255) NOT NULL,
  additional_price DECIMAL(12) NOT NULL,
  total_price DECIMAL(12) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NULL
)