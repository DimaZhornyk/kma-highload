-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE books
(
    `id` int unsigned auto_increment,
    `name` VARCHAR(70),
    `author` varchar(70),
    `created_at`    timestamp NULL,
    `updated_at`    timestamp NULL,
    `deleted_at`    timestamp NULL,
    PRIMARY KEY (`id`)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS `books`;