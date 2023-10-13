create table mesecons(
    poskey varchar(32) primary key not null,
    x int not null,
    y int not null,
    z int not null,
    name varchar(128) not null,
    order_id int not null,
    category varchar(128) not null,
    nodename varchar(128) not null,
    playername varchar(64) not null,
    state varchar(64) not null,
    last_modified bigint not null
);
