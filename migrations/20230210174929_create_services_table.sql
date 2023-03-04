-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS services
(
    id         serial primary key not null,
    name       varchar(255)       not null unique,
    endpoint   varchar(255),
    type       varchar(255)       not null default 'grpc',
    is_active  boolean                     default false,
    deleted_at timestamp(0)                default null,
    created_at timestamp(0)       not null default now(),
    updated_at timestamp(0)       not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS services;
-- +goose StatementEnd