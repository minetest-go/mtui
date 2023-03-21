
create table oauth_app(
    id text(36) primary key not null, -- uuid
    enabled bool not null, -- enabled flag
    created integer not null, -- unix milliseconds
    name text(64) not null, -- application name == clientID
    redirect_urls text not null, -- valid redirect urls, comma-separated
    secret text(64) not null, -- client secret, generated
    allow_privs text not null -- privs to allow onto this application, comman separatedm, empty string means _all_
);
