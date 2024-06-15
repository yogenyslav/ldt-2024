-- +goose Up
-- +goose StatementBegin
alter table chat.favorite_responses
    drop column query_id,
    add column id bigserial primary key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table chat.favorite_responses
    add column query_id bigint primary key generated always as identity,
    drop column id;
-- +goose StatementEnd
