
-- Выданный пропуск
CREATE TABLE issued(
    id int NOT NULL AUTO_INCREMENT COMMENT 'ID выданного пропуска',
    created_at datetime NOT NULL COMMENT 'Дата создания',
    created_by int NOT NULL COMMENT 'Пользователь создания',

    company_inn varchar(255) default null,
    company_ogrn varchar(255) default null,
    company_name text default null,
    company_fio text default null,
    company_car varchar(50) not null,

    legal_basement text default null,
    pass_number varchar(255) not null COMMENT 'Номер пропуска',

    district varchar(255) not null COMMENT 'Наименование МО',

    pass_type int not null COMMENT 'Вид пропуска', -- Краснодар: 1 (желтый); Краснодарский край: 2 (красный)
    issued_at datetime DEFAULT NULL COMMENT 'Дата выдачи',
    registry_number varchar(100) DEFAULT NULL COMMENT 'Номер в реестре округа',
    shipping int not null COMMENT 'Способ направления заявителю', -- 1 - электронно, 2 - нарочно

    file_id int default null,

    PRIMARY KEY (id),
    KEY issued_inn_idx (company_inn),
    KEY issued_ogrn_idx (company_ogrn),
    KEY issued_pass_number_idx (pass_number),
    KEY issued_car_number_idx (company_car),
    KEY issued_registry_number_idx (registry_number),

    CONSTRAINT issued_created_fk FOREIGN KEY (created_by) REFERENCES users (id)

) ENGINE=InnoDB COMMENT 'Выданные пропуска';
