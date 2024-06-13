-- +goose Up
-- +goose StatementBegin
create index session_is_deleted_idx on chat.session (is_deleted) where is_deleted=false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index chat.session_is_deleted_idx;
-- +goose StatementEnd
