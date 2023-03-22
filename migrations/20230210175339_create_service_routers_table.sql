-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS routes
(
    id            serial primary key not null,

    from_method   varchar(255),
    from_address  varchar(255)       not null unique,

    instance      varchar(255),
    proto_service varchar(255),
    proto_method  varchar(255),
    access_role   int                         default 0,

    is_active     boolean                     default true,
    created_at    timestamp(0)       not null default now(),
    updated_at    timestamp(0)       not null default now(),
    deleted_at    timestamp(0)                default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS routes;
-- +goose StatementEnd
