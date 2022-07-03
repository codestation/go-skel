create table test_users
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    name        text        not null,
    external_id uuid        not null,
    primary key (id),
    unique (name)
);

insert into test_users (created_at, updated_at, name, external_id)
values (now(), now(), 'John Doe', '00000000-0000-0000-0000-000000000000'::uuid);

insert into profiles (id, created_at, updated_at, external_id, first_name, last_name, user_token)
    overriding system value
values (1, now(), now(), gen_random_uuid(), 'John', 'Doe', '1234')
on conflict do nothing;
