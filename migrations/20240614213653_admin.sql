-- +goose Up
-- +goose StatementBegin
create schema adm;
create table adm.organization (
    id bigserial primary key ,
    username text unique not null,
    title text unique not null,
    s3_bucket text unique not null,
    created_at timestamp not null default current_timestamp
);
create index organization_username_idx on adm.organization using hash(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table adm.organization;
drop schema adm;
-- +goose StatementEnd
