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

CREATE TABLE IF NOT EXISTS `banners` (
    `id`         BIGINT NOT NULL AUTO_INCREMENT,
    `title`      VARCHAR(256) NOT NULL,
    `body`       TEXT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `active`     BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (`id`)
);

INSERT INTO `banners` (title, body) VALUES
    ("moscow aged", "Московское долголение"),
    ("gold hands", "Клуб очумелые ручки"),
    ("extereme vocals", "Экстрим вокал за неделю"),
    ("ertreme racing", "Курсы экстремального вождения"),
    ("numismatics", "Клуб нумизматов"),
    ("stand up", "Стендап по выходным"),
    ("build complete house", "Дома под ключ"),
    ("flat rent", "Арена квартир"),
    ("delivery food", "Доставка еды"),
    ("delivery items", "Доставка товаров из магазина"),
    ("ship traveling", "Прогулки на теплоходе"),
    ("hotel booking", "Самые выгодные цены на отели"),
    ("car pricing", "Оценка автомобилей"),
    ("places rent", "Аренда торговых площадей"),
    ("sea fishing", "Морская рыбалка"),
    ("diving", "Дайвинг"),
    ("aikido", "Айкидо"),
    ("dance club", "Танцевальная школа"),
    ("dance pair", "Парные танцы"),
    ("holidays agency", "Организация праздников и мероприятий"),
    ("guitar lessons", "Уроки игры на гитаре"),
    ("online courses", "Онлайн обучение"),
    ("english courses", "Обучение английскому"),
    ("group tourism", "Групповые туры"),
    ("farm group", "Клуб садоводов и фермеров"),
    ("gym", "Спортзал с бассеной"),
    ("volunteers", "Волонтёрство"),
    ("concerts", "Лучшие концерты")
;

CREATE TABLE IF NOT EXISTS `slots` (
    `id`         BIGINT NOT NULL AUTO_INCREMENT,
    `title`      VARCHAR(256) NOT NULL,
    `body`       TEXT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `active`     BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (`id`)
);

INSERT INTO `slots` (title, body) VALUES
    ("Main slot", "Слот вверху главного экрана"),
    ("Top slot", "Слот под гланым слотом"),
    ("Right slot", "Слот справа"),
    ("Left slot", "Слот слева"),
    ("Popup slot", "Слот в сплывающем окне")
;

CREATE TABLE IF NOT EXISTS `user_groups` (
    `id`         BIGINT NOT NULL AUTO_INCREMENT,
    `title`      VARCHAR(256) NOT NULL,
    `body`       TEXT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `active`     BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (`id`)
    );

INSERT INTO `user_groups` (title, body) VALUES
    ("women_0_18", "Девочки"),
    ("women_18_30", "Девушки"),
    ("women_30_", "Женщины старше 30"),
    ("women_50_", "Женщины старше 50"),
    ("women_70_", "Женщины старше 70"),
    ("men_0_18", "Мальчики"),
    ("men_18_30", "Юноши"),
    ("men_30_", "Мужчины старше 30"),
    ("men_50_", "Мужчины старше 50"),
    ("men_70_", "Мужчины старше 70")
;

CREATE TABLE IF NOT EXISTS `relations_banner_slot` (
    `id`         BIGINT NOT NULL AUTO_INCREMENT,
    `banner_id`  INTEGER,
    `slot_id`    INTEGER,
    PRIMARY KEY (`id`)
);