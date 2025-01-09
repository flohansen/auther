create table if not exists users (
    username varchar(100) primary key,
    password_hash bytea not null,
    created_at timestamp not null,
    last_modified_at timestamp not null
);
