-- +goose Up
create table game_platforms (
    game_properties_id int NOT NULL,
    constraint game_platforms_game_properties_fk
        FOREIGN KEY (game_properties_id)
        references game_properties(id),
    
    platform_id int NOT NULL,
    constraint game_platforms_platform_fk
        FOREIGN KEY (platform_id)
        references platform(id),
    
    PRIMARY KEY (game_properties_id, platform_id)
);
-- +goose Down
drop table if exists game_platforms;
