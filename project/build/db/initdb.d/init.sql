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

CREATE TABLE IF NOT EXISTS banners (
    id         BIGINT NOT NULL AUTO_INCREMENT,
    title      VARCHAR(256) NOT NULL,
    body       TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active     BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (id)
);

INSERT INTO banners (title, body) VALUES
    ('moscow aged', 'Московское долголение'),
    ('gold hands', 'Клуб очумелые ручки'),
    ('extereme vocals', 'Экстрим вокал за неделю'),
    ('ertreme racing', 'Курсы экстремального вождения'),
    ('numismatics', 'Клуб нумизматов'),
    ('stand up', 'Стендап по выходным'),
    ('build complete house', 'Дома под ключ'),
    ('flat rent', 'Арена квартир'),
    ('delivery food', 'Доставка еды'),
    ('delivery items', 'Доставка товаров из магазина'),
    ('ship traveling', 'Прогулки на теплоходе'),
    ('hotel booking', 'Самые выгодные цены на отели'),
    ('car pricing', 'Оценка автомобилей'),
    ('places rent', 'Аренда торговых площадей'),
    ('sea fishing', 'Морская рыбалка'),
    ('diving', 'Дайвинг'),
    ('aikido', 'Айкидо'),
    ('dance club', 'Танцевальная школа'),
    ('dance pair', 'Парные танцы'),
    ('holidays agency', 'Организация праздников и мероприятий'),
    ('guitar lessons', 'Уроки игры на гитаре'),
    ('online courses', 'Онлайн обучение'),
    ('english courses', 'Обучение английскому'),
    ('group tourism', 'Групповые туры'),
    ('farm group', 'Клуб садоводов и фермеров'),
    ('gym', 'Спортзал с бассеной'),
    ('volunteers', 'Волонтёрство'),
    ('concerts', 'Лучшие концерты')
;

CREATE TABLE IF NOT EXISTS slots (
    id         BIGINT NOT NULL AUTO_INCREMENT,
    title      VARCHAR(256) NOT NULL,
    body       TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active     BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (id)
);

INSERT INTO slots (title, body) VALUES
    ('Main slot', 'Слот вверху главного экрана'),
    ('Top slot', 'Слот под гланым слотом'),
    ('Right slot', 'Слот справа'),
    ('Left slot', 'Слот слева'),
    ('Popup slot', 'Слот в сплывающем окне')
;

CREATE TABLE IF NOT EXISTS user_groups (
    id         BIGINT NOT NULL AUTO_INCREMENT,
    title      VARCHAR(256) NOT NULL,
    body       TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    active     BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (id)
);

INSERT INTO user_groups (title, body) VALUES
    ('women_0_18', 'Девочки'),
    ('women_18_30', 'Девушки'),
    ('women_30_', 'Женщины старше 30'),
    ('women_50_', 'Женщины старше 50'),
    ('women_70_', 'Женщины старше 70'),
    ('men_0_18', 'Мальчики'),
    ('men_18_30', 'Юноши'),
    ('men_30_', 'Мужчины старше 30'),
    ('men_50_', 'Мужчины старше 50'),
    ('men_70_', 'Мужчины старше 70')
;

CREATE TABLE IF NOT EXISTS relations_banner_slot (
    id         BIGINT NOT NULL AUTO_INCREMENT,
    banner_id  INTEGER,
    slot_id    INTEGER,
    PRIMARY KEY (id),
    UNIQUE(banner_id, slot_id)
);

INSERT INTO relations_banner_slot (banner_id, slot_id) VALUES
    (1, 1),
    (1, 7),
    (1, 8),
    (1, 9)
;

CREATE TABLE IF NOT EXISTS events (
    id         BIGINT NOT NULL AUTO_INCREMENT,
    type       VARCHAR(256) NOT NULL,
    banner_id  INTEGER,
    slot_id    INTEGER,
    group_id   INTEGER,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

/**
  показываем баннеры 1 раз каждой группе (инициализация - дерганье за каждую ручку)
 */
INSERT INTO `events` (type, banner_id, slot_id, group_id) VALUES
    /** первый слот*/
    ('show', 1, 1, 1), ('show', 3, 1, 1), ('show', 5, 1, 1),
    ('show', 1, 1, 2), ('show', 3, 1, 2), ('show', 5, 1, 2),
    ('show', 1, 1, 3), ('show', 3, 1, 3), ('show', 5, 1, 3),
    ('show', 1, 1, 4), ('show', 3, 1, 4), ('show', 5, 1, 4),
    ('show', 1, 1, 5), ('show', 3, 1, 5), ('show', 5, 1, 5),
    ('show', 1, 1, 6), ('show', 3, 1, 6), ('show', 5, 1, 6),
    ('show', 1, 1, 7), ('show', 3, 1, 7), ('show', 5, 1, 7),
    ('show', 1, 1, 8), ('show', 3, 1, 8), ('show', 5, 1, 8),
    ('show', 1, 1, 9), ('show', 3, 1, 9), ('show', 5, 1, 9),
    ('show', 1, 1, 10), ('show', 3, 1, 10), ('show', 5, 1, 10),

    /** второй слот*/
    ('show', 2, 2, 1), ('show', 4, 2, 1), ('show', 6, 2, 1),
    ('show', 2, 2, 2), ('show', 4, 2, 2), ('show', 6, 2, 2),
    ('show', 2, 2, 3), ('show', 4, 2, 3), ('show', 6, 2, 3),
    ('show', 2, 2, 4), ('show', 4, 2, 4), ('show', 6, 2, 4),
    ('show', 2, 2, 5), ('show', 4, 2, 5), ('show', 6, 2, 5),
    ('show', 2, 2, 6), ('show', 4, 2, 6), ('show', 6, 2, 6),
    ('show', 2, 2, 7), ('show', 4, 2, 7), ('show', 6, 2, 7),
    ('show', 2, 2, 8), ('show', 4, 2, 8), ('show', 6, 2, 8),
    ('show', 2, 2, 9), ('show', 4, 2, 9), ('show', 6, 2, 9),
    ('show', 2, 2, 10), ('show', 4, 2, 10), ('show', 6, 2, 10),

    /** третий слот*/
    ('show', 7, 3, 1), ('show', 9, 3, 1), ('show', 11, 3, 1),
    ('show', 7, 3, 2), ('show', 9, 3, 2), ('show', 11, 3, 2),
    ('show', 7, 3, 3), ('show', 9, 3, 3), ('show', 11, 3, 3),
    ('show', 7, 3, 4), ('show', 9, 3, 4), ('show', 11, 3, 4),
    ('show', 7, 3, 5), ('show', 9, 3, 5), ('show', 11, 3, 5),
    ('show', 7, 3, 6), ('show', 9, 3, 6), ('show', 11, 3, 6),
    ('show', 7, 3, 7), ('show', 9, 3, 7), ('show', 11, 3, 7),
    ('show', 7, 3, 8), ('show', 9, 3, 8), ('show', 11, 3, 8),
    ('show', 7, 3, 9), ('show', 9, 3, 9), ('show', 11, 3, 9),
    ('show', 7, 3, 10), ('show', 9, 3, 10), ('show', 11, 3, 10),

    /** четвёртый слот*/
    ('show', 8, 4, 1), ('show', 10, 4, 1), ('show', 12, 4, 1),
    ('show', 8, 4, 2), ('show', 10, 4, 2), ('show', 12, 4, 2),
    ('show', 8, 4, 3), ('show', 10, 4, 3), ('show', 12, 4, 3),
    ('show', 8, 4, 4), ('show', 10, 4, 4), ('show', 12, 4, 4),
    ('show', 8, 4, 5), ('show', 10, 4, 5), ('show', 12, 4, 5),
    ('show', 8, 4, 6), ('show', 10, 4, 6), ('show', 12, 4, 6),
    ('show', 8, 4, 7), ('show', 10, 4, 7), ('show', 12, 4, 7),
    ('show', 8, 4, 8), ('show', 10, 4, 8), ('show', 12, 4, 8),
    ('show', 8, 4, 9), ('show', 10, 4, 9), ('show', 12, 4, 9),
    ('show', 8, 4, 10), ('show', 10, 4, 10), ('show', 12, 4, 10),

    /** пятый слот*/
    ('show', 13, 5, 1), ('show', 14, 5, 1), ('show', 15, 5, 1),
    ('show', 13, 5, 2), ('show', 14, 5, 2), ('show', 15, 5, 2),
    ('show', 13, 5, 3), ('show', 14, 5, 3), ('show', 15, 5, 3),
    ('show', 13, 5, 4), ('show', 14, 5, 4), ('show', 15, 5, 4),
    ('show', 13, 5, 5), ('show', 14, 5, 5), ('show', 15, 5, 5),
    ('show', 13, 5, 6), ('show', 14, 5, 6), ('show', 15, 5, 6),
    ('show', 13, 5, 7), ('show', 14, 5, 7), ('show', 15, 5, 7),
    ('show', 13, 5, 8), ('show', 14, 5, 8), ('show', 15, 5, 8),
    ('show', 13, 5, 9), ('show', 14, 5, 9), ('show', 15, 5, 9),
    ('show', 13, 5, 10), ('show', 14, 5, 10), ('show', 15, 5, 10)
;