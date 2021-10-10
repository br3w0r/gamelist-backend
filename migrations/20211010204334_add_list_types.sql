-- +goose Up
insert into list_type(id, name) values
    (1, 'Played'),
    (2, 'Playing'),
    (3, 'Want to play');
-- +goose Down
delete from list_type
    where id in (1, 2, 3);
