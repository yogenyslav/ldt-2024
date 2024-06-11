-- +goose Up
-- +goose StatementBegin
alter table chat.query
    add column period text not null default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table chat.query
    drop column period;
-- +goose StatementEnd
