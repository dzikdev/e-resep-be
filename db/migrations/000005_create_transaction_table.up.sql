CREATE TABLE IF NOT EXISTS transaction (
  id SERIAL NOT NULL PRIMARY KEY,
  status VARCHAR(255) NOT NULL,
  additional_price DECIMAL(12) NOT NULL,
  total_price DECIMAL(12) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)