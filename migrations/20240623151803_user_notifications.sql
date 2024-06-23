-- +goose Up
-- +goose StatementBegin
create table chat.notification (
    email text,
    first_name text,
    last_name text,
    organization_id bigint
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat.notification;
-- +goose StatementEnd
