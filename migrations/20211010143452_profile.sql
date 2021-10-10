-- +goose Up
create table profile (
    id SERIAL PRIMARY KEY,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp,
    nickname varchar(20) UNIQUE NOT NULL,
    description varchar(120),
    games_listed int DEFAULT 0 NOT NULL,
    email varchar(256) NOT NULL,
    password text NOT NULL
);
-- +goose Down
drop table if exists profile;
