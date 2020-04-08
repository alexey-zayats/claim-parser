
DROP TABLE IF EXISTS `passes`;
DROP TABLE IF EXISTS `files`;
DROP TABLE IF EXISTS `bids`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `email` varchar(100) NOT NULL COMMENT 'E-mail',
    `password_hash` varchar(255) NOT NULL COMMENT 'Пароль',
    `password_reset_token` varchar(50) DEFAULT NULL COMMENT 'Токен восстановления пароля',
    `name` varchar(50) NOT NULL COMMENT 'Имя',
    `district` int DEFAULT NULL COMMENT 'Район/округ',
    `pass_type` int DEFAULT NULL COMMENT 'Тип пропуска',
    `created_at` datetime NOT NULL COMMENT 'Дата создания',
    `updated_at` datetime NOT NULL COMMENT 'Дата обновления',
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`),
    UNIQUE KEY `password_reset_token` (`password_reset_token`)
) ENGINE=InnoDB AUTO_INCREMENT=2;

INSERT INTO `users` VALUES (1,'dev@quartex.ru','$2y$13$EZfzE4wEDo.X4qKWOsRk1uJbkHFGZ/z1McRsR94.5Dp82TfrRiN2G',NULL,'QDev',1,1,'2020-04-05 16:25:31','2020-04-05 16:25:31');

CREATE TABLE `files` (
    `id` int NOT NULL AUTO_INCREMENT,
    `filepath` varchar(512) NOT NULL,
    `status` smallint DEFAULT NULL,
    `log` text DEFAULT NULL,
    `source` text DEFAULT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_files_status` (`status`)
) ENGINE=InnoDB COMMENT 'Файлы заявок';

CREATE TABLE `bids` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT 'Идентификатор',
    `file_id` int DEFAULT NULL COMMENT 'ID файла завки',
    `status` int DEFAULT NULL,
    `workflow_status` int DEFAULT NULL,
    `code` text DEFAULT NULL,
    `district` int NOT NULL COMMENT 'Муниципальный округ',
    `type` smallint NOT NULL COMMENT 'Тип пропуска',
    `created_at` datetime NOT NULL COMMENT 'Дата создания',
    `created_by` int NOT NULL COMMENT 'Пользователь создания',
    `user_id` int DEFAULT NULL COMMENT 'Оператор',
    `source` text DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_bids_wstatus` (`workflow_status`)
) ENGINE=InnoDB COMMENT 'Заявки';

CREATE TABLE `passes` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `company_branch` text COMMENT 'Вид деятельности',
    `company_okved` varchar(255) DEFAULT NULL COMMENT 'Код ОКВЭД',
    `company_inn` varchar(100) NOT NULL COMMENT 'ИНН',
    `company_name` text NOT NULL COMMENT 'Название',
    `company_address` text DEFAULT NULL COMMENT 'Адрес',
    `company_ceo_phone` text DEFAULT NULL COMMENT 'Телефон директора',
    `company_ceo_email` text NOT NULL COMMENT 'E-mail директора',
    `company_lastname` varchar(100) NOT NULL COMMENT 'Фамилия директора',
    `company_firstname` varchar(100) NOT NULL COMMENT 'Имя директора',
    `company_patrname` varchar(100) NOT NULL COMMENT 'Отчество директора',
    `employee_lastname` varchar(100) NOT NULL COMMENT 'Фамилия гражданина',
    `employee_firstname` varchar(100) NOT NULL COMMENT 'Имя гражданина',
    `employee_patrname` varchar(100) NOT NULL COMMENT 'Отчество гражданина',
    `employee_car` varchar(20) DEFAULT NULL COMMENT 'Номер автомобиля',
    `employee_agree` smallint NOT NULL COMMENT 'Согласие обработки',
    `employee_confirm` smallint NOT NULL COMMENT 'Подтверждение данных',
    `source` smallint NOT NULL COMMENT 'Источник загрузки',
    `district` int NOT NULL COMMENT 'Муниципальный округ',
    `type` smallint NOT NULL COMMENT 'Тип пропуска',
    `number` varchar(50) DEFAULT NULL COMMENT 'Номер пропуска',
    `alighner_post` varchar(50) DEFAULT NULL COMMENT 'Должность согласователя',
    `alighner_name` varchar(100) DEFAULT NULL COMMENT 'ФИО согласователя',
    `send_type` int DEFAULT NULL COMMENT 'Способ направления',
    `status` smallint NOT NULL COMMENT 'Статус',
    `file_id` int DEFAULT NULL COMMENT 'Файл загрузки',
    `log` text COMMENT 'Журнал обработки',
    `created_at` datetime NOT NULL COMMENT 'Дата создания',
    `created_by` int NOT NULL COMMENT 'Пользователь создания',
    `bid_id` int DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_company_inn` (`company_inn`),
    KEY `idx_district` (`district`),
    KEY `idx_employee_lastname` (`employee_lastname`),
    KEY `fk_passes_created` (`created_by`),
    KEY `idx_passes_district` (`district`),
    KEY `idx_passes_bid` (`bid_id`),
    CONSTRAINT `fk_passes_created` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB COMMENT='Пропуски';


