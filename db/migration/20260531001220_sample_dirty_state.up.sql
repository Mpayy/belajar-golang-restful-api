create table if not exists correct (
    id int auto_increment,
    name varchar(255) not null,
    primary key (id)
) engine=InnoDB;

create table if not exists wrong (
    id int auto_increment,
    name varchar(255) not null,
    primary key (id)
) engine=InnoDB;