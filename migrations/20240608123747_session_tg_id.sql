-- +goose Up
-- +goose StatementBegin
alter table chat.session
    add column tg_id bigint not null default 0;
create index session_tg_id on chat.session using hash(tg_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index chat.session_tg_id;
alter table chat.session
    drop column tg_id;
-- +goose StatementEnd
