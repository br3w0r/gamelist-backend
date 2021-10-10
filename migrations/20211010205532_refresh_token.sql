-- +goose Up
create table refresh_token (
    id SERIAL PRIMARY KEY,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp,
    profile_id int NOT NULL,
    constraint refresh_token_profile_fk
        FOREIGN KEY (profile_id)
        references profile(id),
    
    token varchar(256) UNIQUE NOT NULL
);
-- +goose Down
drop table if exists refresh_token;
