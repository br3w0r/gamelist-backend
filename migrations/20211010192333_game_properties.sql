-- +goose Up
create table game_properties (
    id SERIAL PRIMARY KEY,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp,
    name varchar(256) UNIQUE NOT NULL,
    image_url varchar(2000) NOT NULL,
    year_released smallint NOT NULL
);
-- +goose Down
drop table if exists game_properties;
