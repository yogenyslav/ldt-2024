-- +goose Up
-- +goose StatementBegin
create table chat.query(
    id bigserial primary key,
    prompt text not null default '',
    command text not null default '',
    username text not null,
    session_id uuid not null,
    created_at timestamp not null default current_timestamp
);
create index query_username_idx on chat.query using hash(username);
create index query_session_id_idx on chat.query using hash(session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat.query;
-- +goose StatementEnd
