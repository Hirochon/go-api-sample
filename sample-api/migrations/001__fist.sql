-- +migrate Up
create table `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT unique,
    PRIMARY KEY (`id`)
);

-- +migrate Down
DROP TABLE IF EXISTS user;