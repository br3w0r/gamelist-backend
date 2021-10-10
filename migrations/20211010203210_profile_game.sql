-- +goose Up
create table profile_game (
    profile_id int NOT NULL,
    constraint profile_game_profile_pk
        FOREIGN KEY (profile_id)
        references profile(id),

    game_id int NOT NULL,
    constraint profile_game_game_properties_pk
        FOREIGN KEY (game_id)
        references game_properties(id),

    list_type_id int NOT NULL,
    constraint profile_game_game_list_type_pk
        FOREIGN KEY (list_type_id)
        references list_type(id),

    PRIMARY KEY (profile_id, game_id)
);
-- +goose Down
drop table if exists profile_game;
