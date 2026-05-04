-- +goose Up
-- +goose StatementBegin
ALTER TABLE sales DROP COLUMN IF EXISTS date;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Cannot restore the date column without knowing the original data
-- +goose StatementEnd
