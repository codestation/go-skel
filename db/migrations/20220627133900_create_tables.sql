-- +migrate Up

create table if not exists profiles
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    first_name  text        not null,
    last_name   text        not null,
    email       email       not null,
    primary key (id),
    unique (email),
    check (char_length(first_name) <= 255),
    check (char_length(last_name) <= 255),
    check (char_length(email) <= 254)
);

-- +migrate Down
drop table if exists profiles;
