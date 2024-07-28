CREATE TABLE IF NOT EXISTS medication_request_identifier(
  id SERIAL NOT NULL PRIMARY KEY,
  medication_request_id INT NOT NULL,
  system VARCHAR(255) NOT NULL,
  use VARCHAR(255) NOT NULL,
  value VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)