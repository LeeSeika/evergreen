# drop table if exists user;
# create table user (
#     id bigint(20) not null auto_increment,
#     user_id bigint(20) not null,
#     username varchar(64) collate utf8mb4_general_ci not null,
#     password varchar(64) collate utf8mb4_general_ci not null,
#     email varchar(64) collate utf8mb4_general_ci,
#     gender tinyint(4) not null default 0,
#     create_time timestamp null default current_timestamp,
#     update_time timestamp null default current_timestamp on update current_timestamp,
#     primary key (id),
#     unique key idx_username (username) using btree,
#     unique key idx_user_id (user_id) using btree
# ) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;

drop table if exists community;
create table community (
    id int(11) not null auto_increment,
    community_id int(10) unsigned not null ,
    community_name varchar(128) collate utf8mb4_general_ci not null,
    introduction varchar(256) collate utf8mb4_general_ci not null ,
    create_time timestamp not null default current_timestamp,
    update_time timestamp not null default current_timestamp on update current_timestamp,
    primary key (id),
    unique key idx_community_id (community_id),
    unique key idx_community_name (community_name)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;

insert into community values ('1', '1', 'Go', 'Golang', '2020-01-07 09:10:45', '2020-11-27 19:12:45');
insert into community values ('2', '2', 'Dota2', 'haha', '2020-01-07 09:10:45', '2020-11-27 19:12:45');
insert into community values ('3', '3', 'Java', 'test', '2020-01-07 09:10:45', '2020-11-27 19:12:45');
insert into community values ('4', '4', 'Python', 'hehe', '2020-01-07 09:10:45', '2020-11-27 19:12:45');