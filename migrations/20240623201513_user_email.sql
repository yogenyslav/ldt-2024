-- +goose Up
-- +goose StatementBegin
create table adm.user_email (
    email text not null,
    username text not null
);
insert into adm.user_email (email, username) values ('test@test.com', 'test_admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table adm.user_email;
-- +goose StatementEnd
