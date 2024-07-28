CREATE TABLE IF NOT EXISTS medication_ingredient (
  id SERIAL NOT NULL PRIMARY KEY,
  medication_id INT NOT NULL,
  code VARCHAR(255) NOT NULL,
  display VARCHAR(255) NOT NULL,
  is_active BOOLEAN NOT NULL,
  strength_denominator VARCHAR(255) NOT NULL,
  strength_numerator VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)