-- +goose Up
-- +goose StatementBegin
alter table chat.response
    add column data jsonb not null default '{}',
    add column data_type int8 not null default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table chat.response
    drop column data,
    drop column data_type;
-- +goose StatementEnd
