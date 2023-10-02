
create table chat_log(
    id text(36) primary key not null, -- uuid
    timestamp integer not null, -- unix milliseconds
    channel text(32) not null, -- main
    name text(64) not null, -- from-playername
    message text(512) not null -- message
);

create index chat_log_index on chat_log(timestamp, channel);
