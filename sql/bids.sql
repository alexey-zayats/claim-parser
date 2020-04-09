

CREATE TABLE `bids` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT 'Идентификатор',
    `file_id` int DEFAULT NULL COMMENT 'ID файла завки',
    `workflow_status` int DEFAULT NULL,
    `code` text DEFAULT NULL,
    `district` int NOT NULL COMMENT 'Муниципальный округ',
    `type` smallint NOT NULL COMMENT 'Тип пропуска',
    `created_at` datetime NOT NULL COMMENT 'Дата создания',
    `created_by` int NOT NULL COMMENT 'Пользователь создания',
    `user_id` int DEFAULT NULL COMMENT 'Оператор',
    `source` text DEFAULT NULL,
    `moved_to` int(11) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_bids_wstatus` (`workflow_status`)
) ENGINE=InnoDB COMMENT 'Заявки';


