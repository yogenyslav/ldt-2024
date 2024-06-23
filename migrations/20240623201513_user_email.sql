-- +goose Up
-- +goose StatementBegin
create table adm.user_email (
    email text not null,
    username text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table adm.user_email;
-- +goose StatementEnd
