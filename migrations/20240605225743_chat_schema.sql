-- +goose Up
-- +goose StatementBegin
create schema chat;
alter database dev
    set timezone to 'Europe/Moscow';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema chat;
alter database dev
    set timezone to 'UTC';
-- +goose StatementEnd
