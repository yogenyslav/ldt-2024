-- +goose Up
-- +goose StatementBegin
create table chat.notification (
    email text not null,
    organization_id bigint not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat.notification;
-- +goose StatementEnd
