-- +goose Up
-- +goose StatementBegin
create table chat.response(
    query_id bigint primary key,
    status int4 not null default 1,
    body text not null default '',
    created_at timestamp not null default current_timestamp
);
create index response_status_idx on chat.response (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat.response;
-- +goose StatementEnd
