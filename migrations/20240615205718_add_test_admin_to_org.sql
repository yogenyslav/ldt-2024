-- +goose Up
-- +goose StatementBegin
insert into adm.organization (username, title, s3_bucket)
values ('test_admin', 'default', 'organization-default');
insert into adm.user_organization (username, organization)
values ('test_admin', 'default');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from adm.organization where title = 'default';
delete from adm.user_organization where username = 'test_admin';
-- +goose StatementEnd
