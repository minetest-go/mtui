
alter table log add column nodename varchar(255);

drop index log_index;

create index log_index on log(timestamp, category, event, username, nodename, posx, posy, posz);