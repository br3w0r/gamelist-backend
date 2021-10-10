-- +goose Up
create table genre (
    id SERIAL PRIMARY KEY,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp,
    name varchar(50) UNIQUE NOT NULL
);
-- +goose Down
drop table if exists genre;
