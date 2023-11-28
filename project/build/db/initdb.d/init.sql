CREATE DATABASE IF NOT EXISTS test;
USE test;

CREATE TABLE IF NOT EXISTS `tasks` (
    `id`         BIGINT NOT NULL AUTO_INCREMENT,
    `title`      VARCHAR(256) NOT NULL,
    `body`       TEXT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

INSERT INTO tasks (title, body) VALUES ("task1", "タスク1です"), ("task2", "タスク2です");
