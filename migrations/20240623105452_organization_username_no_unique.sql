-- +goose Up
-- +goose StatementBegin
alter table adm.organization
    drop constraint organization_username_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table adm.organization
    add constraint organization_username_key unique (username);
-- +goose StatementEnd
