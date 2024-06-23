-- +goose Up
-- +goose StatementBegin
alter table adm.organization
    drop column s3_bucket;
alter table adm.user_organization
    add column organization_id bigint not null default 0,
    drop column organization;
create index user_organization_organization_id_idx on adm.user_organization (organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table adm.user_organization
    drop column organization_id,
    add column organization text not null default '';
alter table adm.organization
    add column s3_bucket text not null default '';
-- +goose StatementEnd
