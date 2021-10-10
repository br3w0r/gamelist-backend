-- +goose Up
create table game_genres (
    game_properties_id int NOT NULL,
    constraint game_genres_game_properties_fk
        FOREIGN KEY (game_properties_id)
        references game_properties(id),
    
    genre_id int NOT NULL,
    constraint game_genres_genre_fk
        FOREIGN KEY (genre_id)
        references genre(id),
    
    PRIMARY KEY (game_properties_id, genre_id)
);
-- +goose Down
drop table if exists game_genres;
