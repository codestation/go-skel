create table if not exists test_profiles
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    external_id uuid        not null,
    avatar      text,
    primary key (id)
);

create table if not exists test_users
(
    id          integer generated always as identity,
    created_at  timestamptz not null,
    updated_at  timestamptz not null,
    deleted_at  timestamptz,
    external_id uuid        not null,
    name        text        not null,
    code        integer generated always as identity,
    data1       jsonb       not null,
    data2       jsonb       not null,
    profile_id  integer     not null,
    primary key (id),
    unique (name),
    constraint fk_users_profile foreign key (profile_id) references test_profiles (id)
);

delete from test_profiles;
select setval('test_profiles_id_seq', coalesce((select max(id) from test_profiles), 1), false);

delete from test_users;
select setval('test_users_id_seq', coalesce((select max(id) from test_users), 1), false);

insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000001'::uuid, 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000002'::uuid, 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000003'::uuid, 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000004'::uuid, 'https://example.com/avatar.jpg');
insert into test_profiles (created_at, updated_at, external_id, avatar)
values (now(), now(), '00000000-0000-0000-0000-000000000005'::uuid, 'https://example.com/avatar.jpg');

insert into test_users (created_at, updated_at, name, external_id, data1, data2, profile_id)
values (now(), now(), 'John Doe 1', '00000000-0000-0000-0000-000000000001'::uuid, '{"foo": "one", "bar": 2}', '[{"foo": "one", "bar": 2}]', 1);
insert into test_users (created_at, updated_at, name, external_id, data1, data2, profile_id)
values (now(), now(), 'John Doe 2', '00000000-0000-0000-0000-000000000002'::uuid, '{"foo": "one", "bar": 2}', '[{"foo": "one", "bar": 2}]', 2);
insert into test_users (created_at, updated_at, name, external_id, data1, data2, profile_id)
values (now(), now(), 'John Doe 3', '00000000-0000-0000-0000-000000000003'::uuid, '{"foo": "one", "bar": 2}', '[{"foo": "one", "bar": 2}]', 3);
insert into test_users (created_at, updated_at, name, external_id, data1, data2, profile_id)
values (now(), now(), 'John Doe 4', '00000000-0000-0000-0000-000000000004'::uuid, '{"foo": "one", "bar": 2}', '[{"foo": "one", "bar": 2}]', 4);
insert into test_users (created_at, updated_at, name, external_id, data1, data2, profile_id)
values (now(), now(), 'John Doe 5', '00000000-0000-0000-0000-000000000005'::uuid, '{"foo": "one", "bar": 2}', '[{"foo": "one", "bar": 2}]', 5);

insert into profiles (created_at, updated_at, first_name, last_name, email)
values (now(), now(), 'John', 'Doe', 'john.doe@example.com');
insert into profiles (created_at, updated_at, first_name, last_name, email)
values (now(), now(), 'Jane', 'Doe', 'jane.doe@example.com');
insert into profiles (created_at, updated_at, first_name, last_name, email)
values (now(), now(), 'Jane', 'Smith', 'jane.smith@example.com');
