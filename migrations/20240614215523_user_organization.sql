-- +goose Up
-- +goose StatementBegin
create table adm.user_organization(
    username text not null,
    organization text not null,
    is_deleted bool not null default false,
    created_at timestamp not null default current_timestamp
);
create index user_organization_username_idx on adm.user_organization using hash(username);
create index user_organization_organization_id_idx on adm.user_organization using hash(organization);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table adm.user_organization;
-- +goose StatementEnd
