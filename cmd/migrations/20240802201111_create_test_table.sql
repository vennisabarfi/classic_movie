-- +goose Up
-- +goose StatementBegin
CREATE TABLE test;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test;
-- +goose StatementEnd
