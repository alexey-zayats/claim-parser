
-- Выданный пропуск
CREATE TABLE issued(
    id int NOT NULL AUTO_INCREMENT COMMENT 'ID выданного пропуска',
    district_id int not null COMMENT 'ID округа',
    company_inn varchar(255) default null,
    company_ogrn varchar(255) default null,
    company_name text default null,
    company_fio text default null,
    company_car varchar(15) not null,
    pass_type int not null COMMENT 'вид пропуска',
    pass_number varchar(255) not null COMMENT 'Номер пропуска',
    issued_at datetime DEFAULT NULL COMMENT 'Дата создания',
    created_at datetime NOT NULL COMMENT 'Дата создания',
    created_by int NOT NULL COMMENT 'Пользователь создания',
    PRIMARY KEY (id),
    KEY issued_inn_idx (company_inn),
    KEY issued_ogrn_idx (company_ogrn),
    KEY issued_pass_number_idx (pass_number),
    KEY issued_car_number_idx (company_car),
    CONSTRAINT issued_district_id_fk FOREIGN KEY (district_id) REFERENCES districts (id),
    CONSTRAINT issued_created_fk FOREIGN KEY (created_by) REFERENCES users (id)
) ENGINE=InnoDB COMMENT 'Выданные пропуска';
