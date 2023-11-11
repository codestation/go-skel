delete from profiles;
select setval('profiles_id_seq', COALESCE((SELECT max(id) FROM profiles), 1), false);
