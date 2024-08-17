CREATE TABLE IF NOT EXISTS transaction_detail (
  id SERIAL NOT NULL PRIMARY KEY,
  transaction_id INT NOT NULL,
  medication_id INT NOT NULL,
  medication_name TEXT NOT NULL,
  price DECIMAL(12) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)