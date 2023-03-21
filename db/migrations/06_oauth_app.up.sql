
create table oauth_app(
    id text(36) primary key not null, -- uuid
    enabled bool not null, -- enabled flag
    created integer not null, -- unix milliseconds
    name text(64) not null, -- application name == clientID
    domain text not null, -- valid redirect url
    secret text(64) not null -- client secret, generated
);
