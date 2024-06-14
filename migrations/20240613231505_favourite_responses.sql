-- +goose Up
-- +goose StatementBegin
create table chat.favorite_responses (
    query_id bigint primary key,
    username text not null,
    response jsonb not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    is_deleted bool not null default false
);
create index favorite_username_idx on chat.favorite_responses using hash(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat.favorite_responses;
-- +goose StatementEnd
