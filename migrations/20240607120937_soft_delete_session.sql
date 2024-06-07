-- +goose Up
-- +goose StatementBegin
alter table chat.session
    add column is_deleted bool not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table chat.session
    drop column is_deleted;
-- +goose StatementEnd
