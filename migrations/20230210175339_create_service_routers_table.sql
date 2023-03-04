-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS service_routers
(
    id            serial primary key not null,
    service_id    bigint             not null,

    from_method   varchar(255),
    from_address  varchar(255)       not null unique,

    proto_folder  varchar(255),
    proto_service varchar(255),
    proto_method  varchar(255),

    is_active     boolean                     default true,
    deleted_at    timestamp(0)                default null,
    created_at    timestamp(0)       not null default now(),
    updated_at    timestamp(0)       not null default now(),

    constraint fk_user_id foreign key (service_id) REFERENCES services (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS service_routers;
-- +goose StatementEnd
