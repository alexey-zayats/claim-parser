
CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `email` varchar(100) NOT NULL COMMENT 'E-mail',
    `password_hash` varchar(255) NOT NULL COMMENT 'Пароль',
    `password_reset_token` varchar(50) DEFAULT NULL COMMENT 'Токен восстановления пароля',
    `district` int(11) DEFAULT NULL COMMENT 'Район/округ',
    `created_at` datetime NOT NULL COMMENT 'Дата создания',
    `updated_at` datetime NOT NULL COMMENT 'Дата обновления',
    `lastname` varchar(50) DEFAULT NULL,
    `firstname` varchar(50) DEFAULT NULL,
    `patrname` varchar(50) DEFAULT NULL,
    `post` varchar(50) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`),
    UNIQUE KEY `password_reset_token` (`password_reset_token`)
) ENGINE=InnoDB AUTO_INCREMENT=6;

INSERT INTO `users` VALUES
    (1,'dev@quartex.ru','$2y$13$cxYIxLQLSmi2OuhASQdC8.gEq2I2X8UcnbujfXd98RHuheIOOYgUO',NULL,1,'2020-04-06 17:49:51','2020-04-07 19:24:10','QDev','-','-','-'),
    (3,'dt-dispetcher@quartex.ru','$2y$13$gayzCEDtTP22/V9bRct0veziGhAPnxEr42TOJLiv838M28w43nnjS',NULL,10,'2020-04-08 12:49:24','2020-04-08 12:50:55','Иванов','иван','иванович','главный'),
    (4,'karasun@quartex.ru','$2y$13$U4q7QCzAs.6WKziDqarJN.wg2tBy6tzD2U/e2lg/PrPxorcm1ox46',NULL,2,'2020-04-08 22:33:30','2020-04-08 22:33:30','Карасунский оператор','иван','иваныч','самый главный'),
    (5,'zapad@quartex.ru','$2y$13$jzEt9gQugZaToG0tgM6hQ.P8DhgUVyyMnVb87ReP3upaYZ9BM36mC',NULL,1,'2020-04-08 22:46:28','2020-04-08 22:46:28','Запад','Анна','Ванна','оператор западного');
