
create table log(
    id text(36) primary key not null, -- uuid
    timestamp integer not null, -- unix milliseconds
    category text(32) not null, -- ui, minetest
    event text(32) not null, -- login, join, leave
    username text(64) not null, -- playername/username, SYSTEM
    message text(512) not null,
    ip_address text(128),
    geo_country text(64),
    geo_city text(64),
    posx int,
    posy int,
    posz int,
    attachment text
);

create index log_index on log(timestamp, category, event, username);
