begin;

create extension if not exists pgcrypto;

create or replace function track_times() returns trigger as $$
begin
	new.last_updated := current_timestamp;
	if tg_op='UPDATE' then
		new.created := old.created;
	end if;
	return new;
end;
$$ language plpgsql stable;

create sequence pk_users start with 1000000;

create table if not exists users (
id bigint not null primary key default nextval('pk_users'),
uuid uuid unique not null default gen_random_uuid(),
username varchar(64) not null,
name varchar(256) not null,
created timestamp not null default current_timestamp,
last_updated timestamp not null default current_timestamp
);

create trigger users_lu before update on users 
	for each row execute procedure track_times();

create sequence pk_things start with 1000000;

create table if not exists things (
id bigint not null primary key default nextval('pk_things'),
uuid uuid unique not null default gen_random_uuid(),
owner bigint references users(id) on delete cascade,
name text not null,
value text,
created timestamp not null default current_timestamp,
last_updated timestamp not null default current_timestamp
);

create trigger things_lu before update on things 
	for each row execute procedure track_times();

commit;
