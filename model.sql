-- The name of db is project


create table users (
	user_id serial not null primary key,
	email character varying(384) not null,
	password character varying(60) not null,
	activated_at timestamp with time zone default null,
	created_at timestamp with time zone default current_timestamp
);

create unique index on users (email);

create table activation (
	activation_id serial not null primary key,
	id character varying(36) not null,
	created_at timestamp with time zone default current_timestamp,
	user_id int not null references users (user_id)
);
	
create table courses (
	course_id serial not null primary key,
	name character varying(64) not null,
	price int
);

update users set activated_at = current_timestamp
where user_id = (
	select user_id from activation where id = $1
)
;

--- Trigger1

create or replace function ver() returns trigger language plpgsql as
	$$
		begin
			insert into activation (id, user_id) values (
				uuid_generate_v4()::varchar,
				new.user_id
			);

			return new;
		end;
	$$
;

----

create trigger ver_trigger after insert on users
for each row execute procedure ver() 
;


--- Trigger2

create or replace function done_ver() returns trigger language plpgsql as
	$$
		begin
			if new.activated_at is not null then
				delete from activation where user_id = old.user_id;
			end if;
			return null;
		end;
	$$
;

create trigger done_ver_trigger after update on users
for each row execute procedure done_ver()
;


-------------------

--Mock data for courses

insert into courses (name, price)
	values ('Artificial Intelegence', 4000000.00),
	 ('Data Science', 3000000.00),
	 ('Golang', 2000000.00),
	 ('Java', 1000000.00),
	 ('Front-End Development', 700000.00)
;


--Functions
