-- +migrate Up
create extension if not exists citext;

-- Regex taken from HTML5 email validation spec
-- https://html.spec.whatwg.org/multipage/input.html#e-mail-state-(type=email)
create domain email as citext
    check ( value ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$' );

-- +migrate Down
drop domain email;
drop extension if exists citext;
