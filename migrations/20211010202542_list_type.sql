-- +goose Up
create table list_type (
    id int PRIMARY KEY NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp,
    name varchar(20) UNIQUE NOT NULL
);
-- +goose Down
drop table if exists list_type;
