

-- Справочники ---

-- Виды пропусков
CREATE TABLE pass_kind (
    pass_kind_id int NOT NULL AUTO_INCREMENT COMMENT 'ID типа пропуска',
    name varchar(50) NOT NULL COMMENT 'Латинское написание вида пропуска', -- использовать в коде, ID получить по имени
    title varchar(100) NOT NULL COMMENT 'Русское написание вида пропуска', -- использовать при выводе в интерфейсе
    PRIMARY    KEY (pass_kind_id),
    KEY        pass_kind_name_idx (name)
) ENGINE=InnoDB COMMENT 'Виды пропусков';

-- Источники импорта
CREATE TABLE import_source (
    import_source_id int NOT NULL AUTO_INCREMENT COMMENT 'ID источника импорта',
    name varchar(50) NOT NULL COMMENT 'Латинское написание источника импорта',
    title varchar(100) NOT NULL COMMENT 'Русское написание источника импорта',
    PRIMARY    KEY (import_source_id),
    KEY        import_source_name_idx (name)
) ENGINE=InnoDB COMMENT 'Источники импорта';

-- Статус заявки
CREATE TABLE claim_status (
    claim_status_id int NOT NULL AUTO_INCREMENT COMMENT 'ID статуса заявки',
    name varchar(50) NOT NULL COMMENT 'Латинское написание статуса заявки',
    title varchar(100) NOT NULL COMMENT 'Русское написание статуса заявки',
    PRIMARY    KEY (claim_status_id),
    KEY        claim_status_name_idx (name)
) ENGINE=InnoDB COMMENT 'Статус заявки';

-- Районы/Округа
CREATE TABLE district (
    district_id int NOT NULL AUTO_INCREMENT COMMENT 'ID района',
    name varchar(50) NOT NULL COMMENT 'Латинское написание района',
    title varchar(100) NOT NULL COMMENT 'Русское написание района',
    PRIMARY    KEY (district_id),
    KEY        district_name_idx (name)
) ENGINE=InnoDB COMMENT 'Районы/округа';

-- Пользователи ---
CREATE TABLE users (
    user_id int NOT NULL AUTO_INCREMENT COMMENT 'ID пользователя',
    email varchar(100) NOT NULL COMMENT 'E-mail',
    password_hash varchar(255) NOT NULL COMMENT 'Пароль',
    password_reset_token varchar(50) DEFAULT NULL COMMENT 'Токен восстановления пароля',
    name varchar(50) NOT NULL COMMENT 'Имя',
    district int DEFAULT NULL COMMENT 'Район/округ',
    pass_type int DEFAULT NULL COMMENT 'Тип пропуска',
    created_at datetime NOT NULL COMMENT 'Дата создания',
    updated_at datetime NOT NULL COMMENT 'Дата обновления',
    PRIMARY KEY (user_id),
    UNIQUE KEY email (email),
    UNIQUE KEY password_reset_token (`password_reset_token`)
) ENGINE=InnoDB AUTO_INCREMENT=2 COMMENT 'Пользователи';

INSERT INTO `users` VALUES (1,'dev@quartex.ru','$2y$13$EZfzE4wEDo.X4qKWOsRk1uJbkHFGZ/z1McRsR94.5Dp82TfrRiN2G',NULL,'QDev',1,1,'2020-04-05 16:25:31','2020-04-05 16:25:31');

-- Импорты
CREATE TABLE import (
    import_id  int NOT NULL AUTO_INCREMENT COMMENT 'Идентификатор импорта',
    file       varchar(512) NOT NULL COMMENT 'Файл импорта',
    status     smallint DEFAULT NULL COMMENT 'Статус результата обработки импорта',
    log        text DEFAULT NULL COMMENT 'Лог процесса импорта',
    import_source_id text DEFAULT NULL COMMENT 'Связь с истоником импорта',
    created_at datetime NOT NULL,
    created_by int not null,
    PRIMARY    KEY (import_id),
    KEY        import_status_idx (status),
    CONSTRAINT import_created_at_fk FOREIGN KEY (created_by) REFERENCES users (user_id)
) ENGINE=InnoDB COMMENT 'Импорт файлов заявок';

-- Компании
CREATE TABLE company (
    company_id  int NOT NULL AUTO_INCREMENT COMMENT 'ID компании',
    import_id int DEFAULT NULL COMMENT 'ID файла завки',
    activity text COMMENT 'Вид деятельности',
    okved varchar(255) DEFAULT NULL COMMENT 'Код ОКВЭД',
    inn varchar(100) NOT NULL COMMENT 'ИНН',
    ogrn varchar(100) NOT NULL COMMENT 'ОГРН',
    title text NOT NULL COMMENT 'Название',
    address text DEFAULT NULL COMMENT 'Адрес',
    ceo_phone text DEFAULT NULL COMMENT 'Телефон директора',
    ceo_email text NOT NULL COMMENT 'E-mail директора',
    lastname varchar(100) NOT NULL COMMENT 'Фамилия директора',
    firstname varchar(100) NOT NULL COMMENT 'Имя директора',
    patrname varchar(100) NOT NULL COMMENT 'Отчество директора',
    created_at datetime NOT NULL,
    created_by int not null,
    PRIMARY    KEY (company_id),
    KEY        company_inn_idx (inn),
    KEY        company_ogrn_idx (ogrn),
    CONSTRAINT company_import_id_fk FOREIGN KEY (import_id) REFERENCES import (import_id),
    CONSTRAINT company_created_at_fk FOREIGN KEY (created_by) REFERENCES users (user_id)
) ENGINE=InnoDB COMMENT 'Компании';

-- Заявки
CREATE TABLE claim (
    claim_id int NOT NULL AUTO_INCREMENT COMMENT 'ID заявки ',
    code text DEFAULT NULL,
    created_at datetime NOT NULL COMMENT 'Дата создания',
    created_by int NOT NULL COMMENT 'Пользователь создания',
    district_id int NOT NULL COMMENT 'Муниципальный округ',
    pass_kind_id smallint NOT NULL COMMENT 'Тип пропуска',
    user_id  int DEFAULT NULL COMMENT 'Оператор',
    context  text DEFAULT NULL,
    import_id int DEFAULT NULL COMMENT 'ID файла завки',
    claim_status_id int DEFAULT NULL COMMENT 'Статус процесса обработки заявок',
    agreee text DEFAULT NULL COMMENT 'Согласие обработки',
    confirm text DEFAULT NULL COMMENT 'Подтверждение данных',
    PRIMARY KEY (`claim_id`),
    KEY `claim_claim_status_id_idx` (`claim_status_id`),
    CONSTRAINT claim_import_id_fk FOREIGN KEY (import_id) REFERENCES import (import_id),
    CONSTRAINT claim_claim_status_id_fk FOREIGN KEY (claim_status_id) REFERENCES claim_status (claim_status_id),
    CONSTRAINT claim_pass_kind_id_fk FOREIGN KEY (pass_kind_id) REFERENCES pass_kind (pass_kind_id),
    CONSTRAINT claim_district_id_fk FOREIGN KEY (district_id) REFERENCES district (district_id),
    CONSTRAINT claim_created_at_fk FOREIGN KEY (created_by) REFERENCES users (user_id)
) ENGINE=InnoDB COMMENT 'Заявки';

-- Запрос пропуска
CREATE TABLE pass_request (
    claim_request_id int NOT NULL AUTO_INCREMENT COMMENT 'ID запроса пропуска',
    claim_id int not null,
    created_at datetime NOT NULL COMMENT 'Дата создания',
    created_by int NOT NULL COMMENT 'Пользователь создания',
    employee_lastname varchar(100) NOT NULL COMMENT 'Фамилия гражданина',
    employee_firstname varchar(100) NOT NULL COMMENT 'Имя гражданина',
    employee_patrname varchar(100) NOT NULL COMMENT 'Отчество гражданина',
    employee_car varchar(20) DEFAULT NULL COMMENT 'Номер автомобиля',
    pass_issued_id int DEFAULT NULL,
    PRIMARY KEY (claim_request_id),
    KEY claim_request_employee_car_idx (employee_car),
    CONSTRAINT claim_request_claim_id_fk FOREIGN KEY (claim_id) REFERENCES claim (claim_id),
    CONSTRAINT claim_request_created_at_fk FOREIGN KEY (created_by) REFERENCES users (user_id),
    CONSTRAINT claim_request_pass_issued_id_fk FOREIGN KEY (pass_issued_id) REFERENCES pass_issued (pass_issued_id)
) ENGINE=InnoDB COMMENT 'Запрос пропуска';

-- Выданный пропуск
CREATE TABLE pass_issued(
    pass_issued_id int NOT NULL AUTO_INCREMENT COMMENT 'ID выданного пропуска',
    claim_id int not null,
) ENGINE=InnoDB COMMENT 'Выданные пропуск';