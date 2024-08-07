CREATE TABLE IF NOT EXISTS patient (
  id SERIAL NOT NULL PRIMARY KEY,
  ref_id TEXT NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  phone_number VARCHAR(20) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)