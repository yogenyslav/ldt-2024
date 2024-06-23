-- +goose Up
-- +goose StatementBegin
create table chat.notification (
    email text,
    organization_id bigint
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat.notification;
-- +goose StatementEnd
