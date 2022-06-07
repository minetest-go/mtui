
create table mod(
    id varchar(36) not null, -- uuid
    name varchar(32) not null, -- cdb package
    mod_type int not null, -- game/mod/txp
    source_type int not null, -- cdb/git
    url varchar(256) not null, -- cdb user+package / git url
    version varchar(64) not null, -- cdb version / git branch,tag,commit
    auto_update boolean not null default false
);
