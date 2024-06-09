-- +goose Up
-- +goose StatementBegin
alter table chat.query
    add column type int4 not null default 0,
    add column product text not null default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table chat.query
    drop type,
    drop product;
-- +goose StatementEnd
