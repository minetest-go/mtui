
create table log(
    id text(36) primary key not null, -- uuid
    timestamp integer not null, -- unix milliseconds
    category text(32) not null, -- ui, minetest
    event text(32) not null, -- login, join, leave
    username text(64) not null, -- playername/username, SYSTEM
    message text(512) not null,
    posx int,
    posy int,
    posz int,
    attachment text
);

create index log_index on log(timestamp, category, event, username);
