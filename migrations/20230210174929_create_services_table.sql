-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS services
(
    id         serial primary key not null,

    instance  varchar(255),
    endpoint   varchar(255),

    is_active  boolean                     default true,
    created_at timestamp(0)       not null default now(),
    updated_at timestamp(0)       not null default now(),
    deleted_at timestamp(0)                default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS services;
-- +goose StatementEnd