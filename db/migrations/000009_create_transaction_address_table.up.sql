CREATE TABLE IF NOT EXISTS transaction_address (
  id SERIAL NOT NULL PRIMARY KEY,
  transaction_id INT NOT NULL,
  address TEXT NOT NULL,
  sub_district VARCHAR(255) NOT NULL,
  district VARCHAR(255) NOT NULL,
  city VARCHAR(255) NOT NULL,
  province VARCHAR(255) NOT NULL,
  postal_code VARCHAR(5) NOT NULL,
  coordinates POINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)