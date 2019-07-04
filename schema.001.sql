begin;

create extension if not exists pgcrypto;

create sequence if not exists pk_users start with 1000000;

create table if not exists users (
id bigint not null primary key default nextval('pk_users'),
uuid uuid unique not null default gen_random_uuid(),
username varchar(64) unique not null,
name varchar(256) not null
);

create sequence if not exists pk_things start with 1000000;

create table if not exists things (
id bigint not null primary key default nextval('pk_things'),
uuid uuid unique not null default gen_random_uuid(),
owner bigint references users(id) on delete cascade,
name text not null,
value text
);

insert into users (username,name) values 
	('a@a.com','Mr A'),
	('b@q.com','Mr B'),
	('c@c.com','Mrs C');

insert into things (owner,name,value) values 
	( (select id from users where name='Mr A'), 'Mr A''s house','123000'),
	( (select id from users where name='Mr B'), 'Mr B''s house','45000'),
	( (select id from users where name='Mr B'), 'Mr B''s car','4900'),
	( (select id from users where name='Mrs C'), 'stuff','120');
	
select * from things;

rollback;
