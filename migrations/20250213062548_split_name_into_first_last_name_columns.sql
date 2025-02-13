-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN name TO first_name;
ALTER TABLE users ALTER COLUMN first_name DROP NOT NULL;
ALTER TABLE users ADD COLUMN last_name VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN last_name;
ALTER TABLE users ALTER COLUMN first_name SET NOT NULL;
ALTER TABLE users RENAME COLUMN first_name TO name;
-- +goose StatementEnd
