create table if not exists users
(
    id              serial primary key,
    user_name       varchar(20) not null,
    display_name    varchar(255),
    hashed_password bytea       not null,
    creation_time   timestamp,
    modified_time   timestamp,
    last_login      timestamp
);

create table if not exists articles
(
    id            serial not null primary key,
    title         varchar(255),
    summary       varchar(255),
    content       text,
    published     boolean   default false,
    author        integer,
    creation_time timestamp default CURRENT_TIMESTAMP,
    modified_time timestamp,
    url           varchar(255),
    deleted       boolean   default false
);

create table if not exists pages
(
    id            serial not null primary key,
    title         varchar(255),
    summary       varchar(255),
    content       text,
    published     boolean   default false,
    author        integer,
    creation_time timestamp default CURRENT_TIMESTAMP,
    modified_time timestamp,
    url           varchar(255),
    deleted       boolean   default false
);


create table if not exists tokens
(
    id                varchar(36),
    user_id           integer,
    creation_time     timestamp default current_timestamp,
    last_used         timestamp,
    revoked           boolean   default false,
    revoked_by        integer,
    revoked_timestamp integer,
    foreign key (user_id) references users (id),
    foreign key (revoked_by) references users (id)
);
