-- +goose Up
-- +goose StatementBegin
alter table chat.query
    drop column command,
    add column status int4 not null default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table chat.query
    add column command text not null default '',
    drop column status;
-- +goose StatementEnd
