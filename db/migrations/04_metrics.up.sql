
create table metric_type(
    name text(128) primary key not null,
    type text(32) not null,
    help text(512) not null
);

create table metric(
    timestamp integer not null,
    name text(128) not null references metric_type(name),
    value real not null,
    primary key (timestamp, name)
);
