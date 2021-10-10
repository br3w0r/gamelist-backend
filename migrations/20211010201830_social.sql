-- +goose Up
create table social (
    profile_id int NOT NULL,
    constraint social_profile_fk
        FOREIGN KEY (profile_id)
        references profile(id),

    type_id int NOT NULL,
    constraint social_type_fk
        FOREIGN KEY (type_id)
        references social_type(id),
    
    data varchar(70) NOT NULL,

    PRIMARY KEY (profile_id, type_id)
);
-- +goose Down
drop table if exists social;
