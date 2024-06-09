-- +goose Up
-- +goose StatementBegin
alter table chat.session
    add column tg boolean not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table chat.session
    drop column tg;
-- +goose StatementEnd
