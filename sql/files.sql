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
