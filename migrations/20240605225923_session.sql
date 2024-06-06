-- +goose Up
-- +goose StatementBegin
create table chat.session(
    id uuid primary key,
    username text not null,
    title text not null default '',
    created_at timestamp not null default current_timestamp
);
create index session_username_idx on chat.session using hash(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat.session;
-- +goose StatementEnd
