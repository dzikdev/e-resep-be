CREATE TABLE IF NOT EXISTS medication (
  id SERIAL NOT NULL PRIMARY KEY,
  ref_id TEXT NOT NULL,
  identifier VARCHAR(255) NULL,
  code VARCHAR(255) NULL,
  code_display TEXT NULL,
  form_code VARCHAR(255) NOT NULL,
  form_value VARCHAR(255) NOT NULL,
  amount JSONB NULL,
  status VARCHAR(255) NOT NULL,
  manufacturer VARCHAR(255) NOT NULL,
  extension JSONB NULL,
  batch JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)