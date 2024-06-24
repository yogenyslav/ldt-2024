-- +goose Up
-- +goose StatementBegin
update adm.user_organization
set organization_id = 1
where username = 'test_admin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update adm.user_organization
set organization_id = 0
where username = 'test_admin';
-- +goose StatementEnd
