create table if not exists category (
    id int auto_increment,
    name varchar(255) not null,
    primary key (id)
)engine=InnoDB;