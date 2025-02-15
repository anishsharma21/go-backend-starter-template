-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD CONSTRAINT email_unique UNIQUE (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP CONSTRAINT email_unique;
-- +goose StatementEnd
