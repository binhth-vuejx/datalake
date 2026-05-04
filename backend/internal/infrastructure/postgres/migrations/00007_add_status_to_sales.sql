-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='sales' AND column_name='status') THEN
    ALTER TABLE sales ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'pending';
  END IF;
END $$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_sales_status ON sales(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sales DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
